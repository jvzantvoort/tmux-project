package update

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

func Install(url string) error {
	utils.LogStart()
	defer utils.LogEnd()

	execPath, err := ThisExecutablePath()
	if err != nil {

		return err
	}
	installDir := filepath.Dir(execPath)

	tmpDir, err := os.MkdirTemp("", "tmux-project-update-*")
	if err != nil {
		log.Errorf("failed to create temp directory: %v", err)
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	log.Debugf("tmpDir: %s", tmpDir)

	log.Infof("downloading from %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("failed to download: %v", err)
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("download failed with status: %d", resp.StatusCode)
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	archivePath := filepath.Join(tmpDir, "download.tar.gz")
	out, err := os.Create(archivePath)
	if err != nil {
		log.Errorf("failed to create temp file: %v", err)
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	log.Debugf("archivePath: %s", archivePath)

	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		log.Errorf("failed to save download: %v", err)
		return fmt.Errorf("failed to save download: %w", err)
	}

	log.Infof("unpacking archive")
	if err := unpackTarGz(archivePath, tmpDir); err != nil {
		log.Errorf("failed to unpack archive: %v", err)
		return fmt.Errorf("failed to unpack archive: %w", err)
	}

	log.Infof("installing binaries to %s", installDir)
	if err := installBinaries(tmpDir, installDir); err != nil {
		log.Errorf("failed to install binaries: %v", err)
		return fmt.Errorf("failed to install binaries: %w", err)
	}

	log.Infof("installation completed successfully")
	return nil
}

func unpackTarGz(archivePath, destDir string) error {
	utils.LogStart()
	defer utils.LogEnd()

	file, err := os.Open(archivePath)
	if err != nil {
		utils.Errorf("error %s", err)
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		utils.Errorf("error %s", err)
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			log.Infof("extracting %s", target)
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				utils.Errorf("error %s", err)

				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				utils.Errorf("error %s", err)

				return err
			}
			outFile.Close()
		}
	}

	return nil
}

func selfUpdate(newBinaryPath string) error {
	exePath, err := ThisExecutablePath()
	if err != nil {
		return fmt.Errorf("cannot get executable path: %w", err)
	}

	updateScript := fmt.Sprintf(`
        sleep 1
        mv %s %s
        chmod +x %s
        %s &
    `, newBinaryPath, exePath, exePath, exePath)

	cmd := exec.Command("sh", "-c", updateScript)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start update script: %w", err)
	}

	os.Exit(0) // exit current process to allow replacement
	return nil
}

func safeInstallBinary(srcPath, destPath string) error {
	utils.LogStart()
	defer utils.LogEnd()
	log.Debugf("installing binary from %s to %s", srcPath, destPath)
	updatingSelf := false
	srcPath, _ = ResolvePath(srcPath)
	destPath, _ = ResolvePath(destPath)

	thisExec, err := ThisExecutablePath()
	if err != nil {
		return err
	}
	if destPath == thisExec {
		destPath = destPath + ".new"
		log.Infof("updating the running executable via temporary file %s", destPath)
		updatingSelf = true
	}

	info, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", srcPath, err)
	}

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", srcPath, err)
	}

	if err := os.WriteFile(destPath, data, info.Mode()); err != nil {
		return fmt.Errorf("failed to write %s: %w", destPath, err)
	}

	if updatingSelf {
		if err := selfUpdate(destPath); err != nil {
			return fmt.Errorf("failed to self-update: %w", err)
		}
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
		srcPath := filepath.Join(srcDir, filename)
		destPath := filepath.Join(destDir, filename)

		log.Infof("installing %s", filename)

		err := safeInstallBinary(srcPath, destPath)
		if err != nil {
			return err
		}

	}

	return nil
}
