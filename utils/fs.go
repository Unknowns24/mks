package utils

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	cp "github.com/otiai10/copy"
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
)

// Get this executable path
func GetExecutablePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(exePath), nil
}

// Check file or directory exists
func FileOrDirectoryExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err == nil
	}
	return true
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

func CheckZipIntegrity(zipPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	return nil
}

// Unzip file to destination
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		filePath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), config.FOLDER_PERMISSION); err != nil {
			return err
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		inFile, err := f.Open()
		if err != nil {
			return err
		}
		defer inFile.Close()

		_, err = io.Copy(outFile, inFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func ZipDirectoryContent(outputzipfilepath, directoryPath string) error {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	// Función para agregar archivos y directorios al archivo zip
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Evitar añadir el directorio raíz mismo al archivo zip
		if path == directoryPath {
			return nil
		}

		// Crear la estructura de directorios y archivos en el archivo zip
		relPath, _ := filepath.Rel(directoryPath, path)
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	if err != nil {
		return fmt.Errorf("Error al agregar contenido al zip: %v", err)
	}

	if err := zw.Close(); err != nil {
		return fmt.Errorf("Error al cerrar el archivo zip: %v", err)
	}

	zipFile, err := os.Create(outputzipfilepath)
	if err != nil {
		return fmt.Errorf("Error al crear el archivo zip: %v", err)
	}
	defer zipFile.Close()

	_, err = buf.WriteTo(zipFile)
	if err != nil {
		return fmt.Errorf("Error al escribir el archivo zip: %v", err)
	}
	return nil
}

func SanitizeFileName(input string) string {
	// Reemplazar espacios con guiones bajos
	sanitized := strings.ReplaceAll(input, " ", "_")

	// Usar una expresión regular para eliminar cualquier carácter que no sea alfanumérico, guión o guión bajo
	reg := regexp.MustCompile("[^a-zA-Z0-9-_]+")
	sanitized = reg.ReplaceAllString(sanitized, "_")

	return sanitized
}

// Make temporary directory with random name and return path
func MakeTemporalDirectory() (string, error) {
	//make uuid using google package
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	// create temporal dir path
	tempDirPath := path.Join(global.TemporalsPath, uuid.String())

	// create directory
	err = os.MkdirAll(tempDirPath, config.FOLDER_PERMISSION)

	if err != nil {
		return "", err
	}

	return tempDirPath, nil
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

func ListDirectoriesAndFiles(dirPath string) ([]string, error) {
	var directoriesAndFiles []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		directoriesAndFiles = append(directoriesAndFiles, entry.Name())
	}

	return directoriesAndFiles, nil
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

func ListFiles(dirPath string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
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
