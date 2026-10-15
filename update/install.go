package update

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jvzantvoort/tmux-project/utils"
	log "github.com/sirupsen/logrus"
)

func ResolvePath(path string) (string, error) {
	utils.LogStart()
	defer utils.LogEnd()

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Errorf("failed to resolve absolute path: %v", err)
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}
	return absPath, nil
}

func ThisExecutablePath() (string, error) {
	utils.LogStart()
	defer utils.LogEnd()

	execPath, err := os.Executable()
	if err != nil {
		log.Errorf("failed to get executable path: %v", err)
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	execPath, err = ResolvePath(execPath)
	if err != nil {
		log.Errorf("failed to resolve executable path: %v", err)
		return "", fmt.Errorf("failed to resolve executable path: %w", err)
	}
	return execPath, nil
}

// downloadTimeout limits downloading a release archive or checksum
const downloadTimeout = 5 * time.Minute

// download saves url to path and returns the sha256 of its content
func download(url, path string) (string, error) {
	client := &http.Client{Timeout: downloadTimeout}
	resp, err := client.Get(url) // #nosec G107 - URL is from the GitHub API
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download of %s failed with status: %d", url, resp.StatusCode)
	}

	out, err := os.Create(path) // #nosec G304 - path is within controlled temp dir
	if err != nil {
		return "", fmt.Errorf("failed to create %s: %w", path, err)
	}

	hash := sha256.New()
	_, err = io.Copy(io.MultiWriter(out, hash), resp.Body)
	if cerr := out.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return "", fmt.Errorf("failed to save %s: %w", path, err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// parseChecksum returns the hash from sha256sum output
func parseChecksum(content string) (string, error) {
	fields := strings.Fields(content)
	if len(fields) == 0 {
		return "", fmt.Errorf("empty checksum file")
	}
	sum := strings.ToLower(fields[0])
	if _, err := hex.DecodeString(sum); err != nil || len(sum) != sha256.Size*2 {
		return "", fmt.Errorf("invalid checksum %q", fields[0])
	}
	return sum, nil
}

// Install downloads the release archive from url, verifies it against the
// sha256sum file at checksumURL and installs the binaries it contains into
// installDir.
func Install(url, checksumURL, installDir string) error {
	utils.LogStart()
	defer utils.LogEnd()

	tmpDir, err := os.MkdirTemp("", "tmux-project-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Errorf("failed to remove temp dir: %v", err)
		}
	}()
	log.Debugf("tmpDir: %s", tmpDir)

	checksumPath := filepath.Join(tmpDir, "download.tar.gz.sha256")
	if _, err := download(checksumURL, checksumPath); err != nil {
		return err
	}
	content, err := os.ReadFile(checksumPath) // #nosec G304 - path is within controlled temp dir
	if err != nil {
		return err
	}
	want, err := parseChecksum(string(content))
	if err != nil {
		return err
	}

	log.Infof("downloading from %s", url)
	archivePath := filepath.Join(tmpDir, "download.tar.gz")
	got, err := download(url, archivePath)
	if err != nil {
		return err
	}
	if got != want {
		return fmt.Errorf("checksum mismatch for %s: got %s, want %s", url, got, want)
	}
	log.Infof("checksum verified")

	unpackDir := filepath.Join(tmpDir, "unpacked")
	log.Infof("unpacking archive")
	if err := unpackTarGz(archivePath, unpackDir); err != nil {
		return fmt.Errorf("failed to unpack archive: %w", err)
	}

	log.Infof("installing binaries to %s", installDir)
	if err := installBinaries(unpackDir, installDir); err != nil {
		return fmt.Errorf("failed to install binaries: %w", err)
	}

	log.Infof("installation completed successfully")
	return nil
}

func unpackTarGz(archivePath, destDir string) error {
	utils.LogStart()
	defer utils.LogEnd()

	cleanDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return fmt.Errorf("failed to resolve destination directory: %w", err)
	}

	file, err := os.Open(archivePath) // #nosec G304 - path is within controlled temp dir
	if err != nil {
		utils.Errorf("error %s", err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Errorf("failed to close file: %v", cerr)
		}
	}()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		utils.Errorf("error %s", err)
		return err
	}
	defer func() {
		if cerr := gzr.Close(); cerr != nil {
			log.Errorf("failed to close gzip reader: %v", cerr)
		}
	}()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Prevent Zip Slip vulnerability: validate archive entry path
		cleanName := filepath.Clean(header.Name)

		// Reject absolute paths
		if filepath.IsAbs(cleanName) {
			return fmt.Errorf("invalid archive entry path (absolute): %s", header.Name)
		}

		// Reject paths with .. traversal
		if strings.Contains(cleanName, "..") {
			return fmt.Errorf("invalid archive entry path (traversal): %s", header.Name)
		}

		// Construct and validate target path stays within destination
		target := filepath.Join(cleanDestDir, cleanName)
		cleanTarget := filepath.Clean(target)
		destPrefix := cleanDestDir + string(os.PathSeparator)

		// Ensure target is within destination directory
		if cleanTarget != cleanDestDir && !strings.HasPrefix(cleanTarget, destPrefix) {
			return fmt.Errorf("invalid archive entry path (outside dest): %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(cleanTarget, 0750); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(cleanTarget), 0750); err != nil {
				return err
			}
			log.Infof("extracting %s", cleanTarget)
			outFile, err := os.OpenFile(cleanTarget, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode&0777)) // #nosec G304,G115 - validated path, safe mode mask
			if err != nil {
				utils.Errorf("error %s", err)
				return err
			}
			_, copyErr := io.Copy(outFile, tr) // #nosec G110 - controlled archive source
			if closeErr := outFile.Close(); closeErr != nil {
				utils.Errorf("failed to close output file: %s", closeErr)
			}
			if copyErr != nil {
				utils.Errorf("error %s", copyErr)
				return copyErr
			}
		default:
			// Ignore symlinks, devices, and other special file types for security
			log.Debugf("skipping unsupported file type %v for %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

// installBinary copies srcPath to destPath. The new file is written next
// to destPath and renamed over it, so a running binary, including this one,
// is replaced atomically.
func installBinary(srcPath, destPath string) error {
	utils.LogStart()
	defer utils.LogEnd()
	log.Debugf("installing binary from %s to %s", srcPath, destPath)

	info, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", srcPath, err)
	}

	data, err := os.ReadFile(srcPath) // #nosec G304 - controlled source path
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", srcPath, err)
	}

	tmpPath := destPath + ".new"
	if err := os.WriteFile(tmpPath, data, info.Mode().Perm()); err != nil { // #nosec G703 - destination is the install directory
		return fmt.Errorf("failed to write %s: %w", tmpPath, err)
	}
	if err := os.Rename(tmpPath, destPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to replace %s: %w", destPath, err)
	}
	return nil
}

func installBinaries(srcDir, destDir string) error {
	utils.LogStart()
	defer utils.LogEnd()

	var binaries []string

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Mode()&0111 != 0 && !strings.HasSuffix(path, ".tar.gz") {
			binaries = append(binaries, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(binaries) == 0 {
		return fmt.Errorf("no executable binaries found in archive")
	}

	for _, binary := range binaries {
		filename := filepath.Base(binary)
		destPath := filepath.Join(destDir, filename)

		log.Infof("installing %s", filename)

		err := installBinary(binary, destPath)
		if err != nil {
			return err
		}

	}

	return nil
}
