package utils

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
)

/* ******************************* */
/* ********** VALIDATORS ********* */
/* ******************************* */

// this functions check if file is a "valid" go file (finds package keyword, but not check syntax)
func IsPseudoValidGoFile(filename string) (bool, error) {
	// Lee el contenido del archivo
	content, err := os.ReadFile(filename)
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
func CheckSyntaxGoFile(filename string) (bool, error) {
	fs := token.NewFileSet()
	_, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	return (err == nil), err
}

// this function checks if a certain package name exists in a file (golang file)
func CheckPackageNameInFile(filename, expectedPackageName string) (bool, error) {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filename, nil, parser.PackageClauseOnly)
	if err != nil {
		return false, err
	}
	return node.Name.Name == expectedPackageName, nil
}

// this function checks if a certain function exists in a file (golang file)
func FunctionExistsInFile(filename, functionName string) (bool, error) {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filename, nil, 0)
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
