package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/unknowns24/mks/config"
)

// Create files from templates wich only needs to changes %PACKAGE_NAME%
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

// Read file content and return it
func ReadFile(filePath string) (string, error) {
	// Read template content
	templateContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(templateContent), nil
}

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
