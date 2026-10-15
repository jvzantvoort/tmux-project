package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input  string
		want   [3]int
		wantOk bool
	}{
		{input: "tmux-project-0.18.0", want: [3]int{0, 18, 0}, wantOk: true},
		{input: "v1.2.3", want: [3]int{1, 2, 3}, wantOk: true},
		{input: "0.10.12", want: [3]int{0, 10, 12}, wantOk: true},
		{input: "dev", wantOk: false},
		{input: "v0.0.0-20261015-abcdef", wantOk: false},
		{input: "", wantOk: false},
	}

	for _, tt := range tests {
		got, ok := parseVersion(tt.input)
		if ok != tt.wantOk || (ok && got != tt.want) {
			t.Errorf("parseVersion(%q) = %v, %v; want %v, %v", tt.input, got, ok, tt.want, tt.wantOk)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b [3]int
		want int
	}{
		{a: [3]int{0, 18, 0}, b: [3]int{0, 18, 0}, want: 0},
		{a: [3]int{0, 9, 0}, b: [3]int{0, 18, 0}, want: -1},
		{a: [3]int{0, 18, 1}, b: [3]int{0, 18, 0}, want: 1},
		{a: [3]int{1, 0, 0}, b: [3]int{0, 99, 99}, want: 1},
	}
	for _, tt := range tests {
		if got := compareVersions(tt.a, tt.b); got != tt.want {
			t.Errorf("compareVersions(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestFindAssets(t *testing.T) {
	release := &GitHubRelease{
		TagName: "tmux-project-0.18.0",
		Assets: []GitHubAsset{
			{Name: "tmux-project.0.18.0.linux.amd64.tar.gz", BrowserDownloadURL: "https://example/amd64"},
			{Name: "tmux-project.0.18.0.linux.amd64.tar.gz.sha256", BrowserDownloadURL: "https://example/amd64.sha256"},
			{Name: "tmux-project.0.18.0.linux.arm64.tar.gz", BrowserDownloadURL: "https://example/arm64"},
		},
	}

	archive, checksum, err := findAssets(release, "linux", "amd64")
	if err != nil {
		t.Fatalf("findAssets() error = %v", err)
	}
	if archive.BrowserDownloadURL != "https://example/amd64" || checksum.BrowserDownloadURL != "https://example/amd64.sha256" {
		t.Errorf("findAssets() = %v, %v", archive, checksum)
	}

	if _, _, err := findAssets(release, "linux", "arm64"); err == nil {
		t.Error("findAssets() without checksum expected an error")
	}
	if _, _, err := findAssets(release, "darwin", "amd64"); err == nil {
		t.Error("findAssets() for missing platform expected an error")
	}
}

func TestParseChecksum(t *testing.T) {
	sum := sha256.Sum256([]byte("x"))
	valid := hex.EncodeToString(sum[:])

	if got, err := parseChecksum(valid + "  file.tar.gz\n"); err != nil || got != valid {
		t.Errorf("parseChecksum() = %q, %v; want %q", got, err, valid)
	}
	for _, input := range []string{"", "nothex  file", "abcd  file"} {
		if _, err := parseChecksum(input); err == nil {
			t.Errorf("parseChecksum(%q) expected an error", input)
		}
	}
}

// buildArchive returns a tar.gz holding the given executables
func buildArchive(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)
	for name, content := range files {
		hdr := &tar.Header{Name: name, Mode: 0755, Size: int64(len(content)), Typeflag: tar.TypeReg}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func serve(t *testing.T, archive []byte, checksum string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/archive.tar.gz", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(archive) })
	mux.HandleFunc("/archive.tar.gz.sha256", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(checksum + "  archive.tar.gz\n"))
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

func TestInstall(t *testing.T) {
	archive := buildArchive(t, map[string]string{"./tmux-project": "new binary", "./proj_info": "new info"})
	sum := sha256.Sum256(archive)
	server := serve(t, archive, hex.EncodeToString(sum[:]))

	installDir := t.TempDir()
	existing := filepath.Join(installDir, "tmux-project")
	if err := os.WriteFile(existing, []byte("old binary"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := Install(server.URL+"/archive.tar.gz", server.URL+"/archive.tar.gz.sha256", installDir); err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	for name, want := range map[string]string{"tmux-project": "new binary", "proj_info": "new info"} {
		path := filepath.Join(installDir, name)
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%s not installed: %v", name, err)
		}
		if string(got) != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
		info, _ := os.Stat(path)
		if info.Mode().Perm()&0100 == 0 {
			t.Errorf("%s is not executable: %v", name, info.Mode())
		}
	}

	if leftovers, _ := filepath.Glob(filepath.Join(installDir, "*.new")); len(leftovers) != 0 {
		t.Errorf("temporary files left behind: %v", leftovers)
	}
}

func TestInstallChecksumMismatch(t *testing.T) {
	archive := buildArchive(t, map[string]string{"./tmux-project": "tampered"})
	sum := sha256.Sum256([]byte("something else"))
	server := serve(t, archive, hex.EncodeToString(sum[:]))

	installDir := t.TempDir()
	existing := filepath.Join(installDir, "tmux-project")
	if err := os.WriteFile(existing, []byte("old binary"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := Install(server.URL+"/archive.tar.gz", server.URL+"/archive.tar.gz.sha256", installDir); err == nil {
		t.Fatal("Install() expected a checksum error")
	}

	got, _ := os.ReadFile(existing)
	if string(got) != "old binary" {
		t.Errorf("binary was replaced despite checksum mismatch: %q", got)
	}
}

func TestUnpackTarGzRejectsTraversal(t *testing.T) {
	archive := buildArchive(t, map[string]string{"../evil": "x"})
	path := filepath.Join(t.TempDir(), "a.tar.gz")
	if err := os.WriteFile(path, archive, 0600); err != nil {
		t.Fatal(err)
	}
	if err := unpackTarGz(path, filepath.Join(t.TempDir(), "out")); err == nil {
		t.Error("unpackTarGz() expected an error for a path traversal entry")
	}
}
