package utils

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
)

/* ******************************* */
/* ********** VALIDATORS ********* */
/* ******************************* */

func TempFileWithDummyPlaceholder(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// CHeck place holders in file %%([A-Z_]+)%% Can't be start with _, and length must be at least 1
	reCheck := regexp.MustCompile("%%(_[A-Z_]+|[A-Z_]{0})%%")
	if matches := reCheck.FindAll(content, -1); len(matches) > 0 {
		return "", errors.New("Placeholders in file can't start with _, and length must be at least 1")
	}

	// Reemplazamos todas las ocurrencias de %%([A-Z_]+)%% con el string resultante
	reReplace := regexp.MustCompile("%%([A-Z_]+)%%")
	modifiedContent := reReplace.ReplaceAll(content, []byte("$1"))

	tempDir, err := MakeTempDirectory()
	if err != nil {
		return "", err
	}

	// Creando el archivo temporal con el mismo nombre en el directorio temporal
	tempFilePath := tempDir + "/" + filepath.Base(filePath)
	err = os.WriteFile(tempFilePath, modifiedContent, 0644)
	if err != nil {
		return "", err
	}

	return tempFilePath, nil
}

// this functions check if file is a "valid" go file (finds package keyword, but not check syntax)
func IsPseudoValidGoFile(filePath string) (bool, error) {
	// Lee el contenido del archivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	// Expresión regular para validar el archivo Go
	// ^(?:[\s]*|//.*\n|/\*.*?\*/)*package\s+\w+
	// ^: Comienza al inicio del texto
	// (?:...): Grupo no capturador
	// [\s]*: Cualquier cantidad de espacios en blanco (incluidos saltos de línea)
	// //.*\n: Comentario de línea
	// /\*.*?\*/: Comentario multilínea (no greedy)
	// package\s+\w+: La línea "package XXXXXX"
	pattern := `^(?:[\s]*|//.*\n|/\*.*?\*/)*package\s+\w+`
	regex := regexp.MustCompile(pattern)

	return regex.Match(content), nil
}

// this function checks if is a valid go file
func CheckSyntaxGoFile(filePath string) (bool, error) {
	tempFilePath, err := TempFileWithDummyPlaceholder(filePath)
	if err != nil {
		return false, err
	}

	fs := token.NewFileSet()
	_, err = parser.ParseFile(fs, tempFilePath, nil, parser.AllErrors)
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
	return false, nil
}
