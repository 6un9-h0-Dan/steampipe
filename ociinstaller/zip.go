package ociinstaller

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"
)

func ungzip(sourceFile string, destDir string) (string, error) {
	r, err := os.Open(sourceFile)
	if err != nil {
		return "", err
	}

	uncompressedStream, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}

	destFile := filepath.Join(destDir, uncompressedStream.Name)
	outFile, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(outFile, uncompressedStream); err != nil {
		return "", err
	}

	outFile.Close()
	if err := uncompressedStream.Close(); err != nil {
		return "", err
	}

	return destFile, nil
}

func unzip(src, dst string) ([]string, error) {
	var files []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return files, err
	}
	defer r.Close()

	os.MkdirAll(dst, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dst, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			if _, err = io.Copy(f, rc); err != nil {
				return err
			}
			f.Close()
		}

		return nil
	}

	for _, f := range r.File {
		if err := extractAndWriteFile(f); err != nil {
			return files, err
		}
		files = append(files, f.FileHeader.Name)
	}

	return files, nil
}

func untar(src, dst string) error {
	fReader, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fReader.Close()

	xzReader, err := xz.NewReader(fReader)
	if err != nil {
		return err
	}

	// create the tar reader from XZ reader
	tarReader := tar.NewReader(xzReader)

	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		fmt.Print(".")

		path := filepath.Join(dst, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		ensureParentPath(path, 0755)

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}

		if _, err = io.Copy(file, tarReader); err != nil {
			return err
		}

		file.Close()

	}
	return nil
}

func ensureParentPath(path string, fileMode os.FileMode) error {
	parentPath := filepath.Dir(path)
	_, err := os.Stat(parentPath)
	if os.IsNotExist(err) {
		return os.MkdirAll(parentPath, fileMode)
	}
	return err
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func copyFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	if _, err = io.Copy(outputFile, inputFile); err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}

	// copy over the permissions and modes
	inputStat, _ := os.Stat(sourcePath)
	outputFile.Chmod(inputStat.Mode())

	return nil
}

func copyFileUnlessExists(sourcePath string, destPath string) error {
	if fileExists(destPath) {
		return nil
	}
	return copyFile(sourcePath, destPath)
}

// moves a file within an fs partition. panics if movement is attempted between partitions
// this is done separately to achieve performance benefits of os.Rename over reading and writing content
func moveFileWithinPartition(sourcePath, destPath string) error {
	if err := os.Rename(sourcePath, destPath); err != nil {
		return fmt.Errorf("error moving file: %s", err)
	}
	return nil
}

// moves a folder within an fs partition. panics if movement is attempted between partitions
// this is done separately to achieve performance benefits of os.Rename over reading and writing content
func moveFolderWithinPartition(sourcePath, destPath string) error {
	sourceinfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(destPath, sourceinfo.Mode()); err != nil {
		return err
	}

	directory, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source dir: %s", err)
	}
	directory.Close()

	defer os.RemoveAll(sourcePath)

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		relPath, _ := filepath.Rel(sourcePath, path)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.MkdirAll(filepath.Join(destPath, relPath), info.Mode())
		}
		return moveFileWithinPartition(filepath.Join(sourcePath, relPath), filepath.Join(destPath, relPath))
	})
}
