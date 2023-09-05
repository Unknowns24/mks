package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"
)

// chech if is an url or not (accepts http and https)
func IsUrl(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// chech if is an github url  (package style like github.com/unknowns24/mks)
func IsGithubUrl(url string) bool {
	return strings.HasPrefix(url, "github.com/")
}

// Check file or directory exists
func FileOrDirectoryExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err == nil
	}
	return true
}

// Make directory if not exists (including parents)
func MakeDirectoryIfNotExists(dirPath string, perms os.FileMode) error {
	if perms == 0 {
		perms = 0755
	}

	if !FileOrDirectoryExists(dirPath) {
		return os.MkdirAll(dirPath, perms)
	}
	// return error if directory already exists
	return errors.New("directory already exists")
}

// move file or directory (with contents) to new location
func MoveFileOrDirectory(sourcePath string, destPath string) error {
	if !FileOrDirectoryExists(sourcePath) {
		return errors.New("source file or directory does not exist")
	}
	if !FileOrDirectoryExists(destPath) {
		return os.Rename(sourcePath, destPath)
	}
	return errors.New("destination file or directory already exists")
}

// Copy file or directory (with contents) to new location
func CopyFileOrDirectory(sourcePath string, destPath string) error {
	if !FileOrDirectoryExists(sourcePath) {
		return errors.New("source file or directory does not exist")
	}
	if !FileOrDirectoryExists(destPath) {
		return cp.Copy(sourcePath, destPath)
	}
	return errors.New("destination file or directory already exists")
}

// Delete file or directory (with contents)
func DeleteFileOrDirectory(filePath string) error {
	if !FileOrDirectoryExists(filePath) {
		return errors.New("file or directory does not exist")
	}
	return os.RemoveAll(filePath)
}

// Download file from url to destination
func DownloadFile(url string, destPath string) error {
	// Create file
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Download file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Unzip file to destination
func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
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
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

// Make temporary directory with random name and return path
func MakeTempDirectory() (string, error) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

// Read file content and return it
func ReadFile(filePath string) (string, error) {
	// Read template content
	templateContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(templateContent), nil
}

func ListDirectories(dirPath string) ([]string, error) {
	var directories []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}

	return directories, nil
}

func findClosingBrace(lines []string, startIndex int) int {
	braceCount := 0

	for i := startIndex; i < len(lines); i++ {
		line := lines[i]
		braceCount += strings.Count(line, "{")
		braceCount -= strings.Count(line, "}")

		if braceCount == 0 {
			return i
		}
	}

	return -1
}
