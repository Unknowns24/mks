package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
)

/******************************
* CREATE FILES FROM TEMPLATES *
*******************************/

// Create files from templates wich only needs to changes %%PACKAGE_NAME%%
func CreateFileFromTemplate(templatePath, serviceName, finalPath string) error {
	// Read template content
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Replace placeholders with the actual service name
	fileContent := strings.ReplaceAll(string(templateContent), config.PLACEHOLDER_PACKAGENAME, serviceName)

	// Write the file content in the file
	if err := os.WriteFile(finalPath, []byte(fileContent), os.ModePerm); err != nil {
		return err
	}

	return nil
}

// Create files from templates that has many placeholders to replace
func CreateFileFromTemplateWithCustomReplace(templatePath, finalPath string, replaces map[string]string) error {
	// Read template content
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fileContent := string(templateContent)

	// Replace every placeholders with its value
	for flag, value := range replaces {
		fileContent = strings.ReplaceAll(string(fileContent), flag, value)
	}

	// Write the file content in the file
	if err := os.WriteFile(finalPath, []byte(fileContent), os.ModePerm); err != nil {
		return err
	}

	return nil
}

/****************************************
* READ AND REPLACE PLACEHOLDERS ON FILE *
*****************************************/

// Read file content, replace placeholder with its values in memory and returns modified content
func ReadFileWithCustomReplace(filePath string, replaces map[string]string) (string, error) {
	// Read template content
	templateContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	fileContent := string(templateContent)

	// Replace every placeholders with its value
	for flag, value := range replaces {
		fileContent = strings.ReplaceAll(string(fileContent), flag, value)
	}

	return fileContent, nil
}

/***********************
* EXTENDS FILE CONTENT *
************************/

// Function to extend file content by adding new content at the bottom of the file
func ExtendFile(filePath, newContent string) error {
	// Read the current content of the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	finalContent := string(content)

	// Check if the file extension need more imports
	if strings.Contains(newContent, "import") {
		mainImports := ExtractImports(finalContent)
		newImports := ExtractImports(newContent)

		// Parsing and creting new headers
		finalHeader := ExtractPackageLine(finalContent)
		finalHeader += "\nimport ("

		for _, importLine := range mainImports {
			finalHeader += "\n" + importLine
		}

		for _, importLine := range newImports {
			finalHeader += "\n" + importLine
		}

		finalHeader += "\n)"

		cleanMainContent := RemovePackageAndImports(finalContent)
		cleanNewContent := RemovePackageAndImports(newContent)

		// Creating final Content
		finalContent = fmt.Sprintf("%s\n%s\n\n%s", finalHeader, cleanMainContent, cleanNewContent)
	} else {
		finalContent = fmt.Sprintf("%s\n\n%s", content, newContent)
	}

	// Write the updated content to the env file
	err = ioutil.WriteFile(filePath, []byte(finalContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

/**********************************
* ADD NEW CONFIGS TO CONFIG FILES *
***********************************/

// Function to add a new configuration before the closing of the Config structure in the source file
func AddGoConfigFromString(newConfig string) error {
	// Path to the source file
	filePath := filepath.Join(global.BasePath, config.FOLDER_SRC, config.FOLDER_UTILS, config.FILE_GO_CONFIG)

	// Read the current content of the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var configClosingBraceIndex int

	// Find the line that defines the Config structure
	for i, line := range lines {
		if strings.Contains(line, "Config") && strings.Contains(line, "struct") {
			configClosingBraceIndex = findClosingBrace(lines, i)
			break
		}
	}

	if configClosingBraceIndex == -1 {
		return fmt.Errorf("could not find closing brace for the Config structure")
	}

	// Insert the new configuration before the closing of the Config structure
	newContent := fmt.Sprintf("%s\n%s\n%s", strings.Join(lines[:configClosingBraceIndex], "\n"), newConfig, strings.Join(lines[configClosingBraceIndex:], "\n"))

	// Write the updated content to the file
	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Function to add a new configuration before the closing of the Config structure in the source file
func AddEnvConfigFromString(newConfig string) error {
	// Path to the source file
	envFilePath := filepath.Join(global.BasePath, config.FILE_CONFIG_ENV)
	exampleEnvfilePath := filepath.Join(global.BasePath, config.FILE_CONFIG_ENVEXAMPLE)

	// Read the current content of the file
	content, err := ioutil.ReadFile(envFilePath)
	if err != nil {
		return err
	}

	// Insert the new configuration before the closing of the Config structure
	newContent := fmt.Sprintf("%s\n%s", content, newConfig)

	// Write the updated content to the env file
	err = ioutil.WriteFile(envFilePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	// Write the updated content to the example env file
	err = ioutil.WriteFile(exampleEnvfilePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

/******************************
* EXTRACT OR REMOVE FROM CODE *
*******************************/

// ExtractPackageLine extracts the line containing the "package" declaration from Go code.
func ExtractPackageLine(code string) string {
	// Utiliza una expresión regular para encontrar la línea que comienza con "package" seguida de caracteres opcionales.
	packageRegex := regexp.MustCompile(`(?m)^package.*\n`)
	matches := packageRegex.FindStringSubmatch(code)

	if len(matches) > 0 {
		return matches[0]
	}

	return ""
}

// RemovePackageAndImports removes the package declaration and all import statements from Go code.
func RemovePackageAndImports(code string) string {
	// Remove the line starting with "package" followed by optional characters.
	packageRegex := regexp.MustCompile(`(?m)^package.*\n`)
	code = packageRegex.ReplaceAllString(code, "")

	// Remove all lines starting with "import" followed by optional characters.
	importRegex := regexp.MustCompile(`import\s+\(\s*([\s\S]*?)\s*\)`)
	code = importRegex.ReplaceAllString(code, "")

	return code
}

// ExtractImports extracts import statements from Go code.
func ExtractImports(code string) []string {
	imports := make([]string, 0)
	// The regular expression finds all lines that start with "import," followed by optional characters,
	// and ends with a semicolon or a newline.
	regex := regexp.MustCompile(`import\s+\(\s*([\s\S]*?)\s*\)`)
	matches := regex.FindAllStringSubmatch(code, -1)

	for _, match := range matches {
		imports = append(imports, match[1])
	}

	return imports
}
