package utils

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
)

/* ******************************* */
/* ********** VALIDATORS ********* */
/* ******************************* */

func TempFileWithDummyPlaceholder(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Check place holders in file %%([A-Z_]+)%% Can't be start with _, and length must be at least 1
	reCheck := regexp.MustCompile("%%(_[A-Z_]+|[A-Z_]{0})%%")
	if matches := reCheck.FindAll(content, -1); len(matches) > 0 {
		return "", errors.New("placeholders in file can't start with _, and length must be at least 1")
	}

	// Replace all the occurencies of %%([A-Z_]+)%% with resultant string
	reReplace := regexp.MustCompile("%%([A-Z_]+)%%")
	modifiedContent := reReplace.ReplaceAll(content, []byte("$1"))

	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	newFileName := uuid.String() + "-" + filepath.Base(filePath)

	// Final path of the temp file will be -> tempDir/fileName
	finalTempFilePath := path.Join(global.TemporalsPath, newFileName)

	// Save the file on the temp folder
	err = os.WriteFile(finalTempFilePath, modifiedContent, config.FOLDER_PERMISSION)
	if err != nil {
		return "", err
	}

	return finalTempFilePath, nil
}

// this function checks if is a valid go file
func CheckSyntaxGoFile(filePath string) (bool, error) {
	tempFilePath, err := TempFileWithDummyPlaceholder(filePath)
	if err != nil {
		return false, err
	}

	fs := token.NewFileSet()
	_, err = parser.ParseFile(fs, tempFilePath, nil, parser.AllErrors)

	DeleteFileOrDirectory(tempFilePath)

	return (err == nil), err
}

// this function checks if a certain package name exists in a file (golang file)
func CheckPackageNameInFile(filePath string, expectedPackageName string) (bool, error) {
	tempFilePath, err := TempFileWithDummyPlaceholder(filePath)
	if err != nil {
		return false, err
	}

	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, tempFilePath, nil, parser.PackageClauseOnly)
	if err != nil {
		return false, err
	}

	DeleteFileOrDirectory(tempFilePath)

	return node.Name.Name == expectedPackageName, nil
}

// this function checks if a certain function exists in a file (golang file)
func FunctionExistsInFile(filePath string, functionName string) (bool, error) {
	tempFilePath, err := TempFileWithDummyPlaceholder(filePath)
	if err != nil {
		return false, err
	}

	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, tempFilePath, nil, 0)
	if err != nil {
		return false, err
	}

	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.Name == functionName {
				return true, nil
			}
		}
	}

	DeleteFileOrDirectory(tempFilePath)

	return false, nil
}
