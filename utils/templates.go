package utils

import (
	"fmt"
	"os"
	"path"
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
func CreateFileFromTemplate(templatePath, finalPath string) error {
	// Read template content
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Replace placeholders with the actual application name
	fileContent := strings.ReplaceAll(string(templateContent), config.PLACEHOLDER_PACKAGENAME, global.ApplicationName)

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
	content, err := os.ReadFile(filePath)
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
	err = os.WriteFile(filePath, []byte(finalContent), config.FOLDER_PERMISSION)
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
	content, err := os.ReadFile(filePath)
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
	err = os.WriteFile(filePath, []byte(newContent), config.FOLDER_PERMISSION)
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
	content, err := os.ReadFile(envFilePath)
	if err != nil {
		return err
	}

	// Insert the new configuration before the closing of the Config structure
	newContent := fmt.Sprintf("%s\n%s", content, newConfig)

	// Write the updated content to the env file
	err = os.WriteFile(envFilePath, []byte(newContent), config.FOLDER_PERMISSION)
	if err != nil {
		return err
	}

	// Write the updated content to the example env file
	err = os.WriteFile(exampleEnvfilePath, []byte(newContent), config.FOLDER_PERMISSION)
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

func ImportBaseContent(sourcePath, finalPath string) error {
	filesAndDirs, err := ListDirectoriesAndFiles(sourcePath)
	if err != nil {
		return err
	}

	// Parse every file or dir
	for _, fileOrDir := range filesAndDirs {
		fileOrDirPath := path.Join(sourcePath, fileOrDir)

		// Check if fileOrDir is mks_modules folder
		if fileOrDir == config.FOLDER_MKS_MODULES {
			if global.Verbose {
				fmt.Println("[+] Copying mks_modules folder..")
			}

			folderFinalPath := path.Join(finalPath, fileOrDir)
			err = CopyFileOrDirectory(fileOrDirPath, folderFinalPath)
			if err != nil {
				return err
			}

			continue
		}

		// Get path info
		info, err := os.Stat(fileOrDirPath)
		if err != nil {
			return err
		}

		// Action depends on file type
		if info.IsDir() {
			if global.Verbose {
				fmt.Printf("[+] Creating %s directory..\n", fileOrDir)
			}

			// Create folder on application
			finalFolderPath := filepath.Join(finalPath, fileOrDir)
			if err := os.MkdirAll(finalFolderPath, os.ModePerm); err != nil {
				return err
			}

			// Scan folder content
			ImportBaseContent(fileOrDirPath, finalFolderPath)
		} else {
			isTemplate := false // <- flag to check if use CreateFileFromTemplate or no

			// If file extension is .template change it for .go on final file
			finalFile := fileOrDir
			if strings.HasSuffix(finalFile, config.FILE_EXTENSION_TEMPLATE) {
				finalFile = strings.ReplaceAll(finalFile, config.FILE_EXTENSION_TEMPLATE, config.FILE_EXTENSION_GO)
				isTemplate = true
			}

			if global.Verbose {
				fmt.Printf("[+] Creating %s file..\n", finalFile)
			}

			// Check if is template to crete file from it or just copy the file
			finalFilePath := path.Join(finalPath, finalFile)
			if isTemplate {
				err := CreateFileFromTemplate(fileOrDirPath, finalFilePath)
				if err != nil {
					return err
				}
			} else {
				CopyFileOrDirectory(fileOrDirPath, finalFilePath)
			}
		}
	}

	return nil
}
