package utils

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/validators"
)

type Prompt struct {
	Type        string `json:"type"`
	Prompt      string `json:"prompt"`
	Default     string `json:"default"`
	Validate    string `json:"validate"`
	ExtendFile  string `json:"extendFile"`
	Placeholder string `json:"placeholder"`
}

type PromptsFileFormat struct {
	Prompts []Prompt `json:"prompts"`
}

func extractPromptValidationParamsValues(input string) []string {
	// Regular expression to identify content within parentheses
	re := regexp.MustCompile(`\(([^)]+)\)`)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return nil
	}

	// Extract and split the values that are separated by commas
	values := regexp.MustCompile(`\s*,\s*`).Split(matches[1], -1)

	return values
}

func getRealPromptValidationParamType(paramValue string) interface{} {
	// Check if paramValue is an int
	intValue, err := strconv.Atoi(paramValue)
	if err == nil {
		return intValue
	}

	// Check if paramValue is boolean
	boolValue, err := strconv.ParseBool(paramValue)
	if err == nil {
		return boolValue
	}

	// if paramValue is not boolean or integer return it as string
	return paramValue
}

func noValidationReplacePrompt(prompt Prompt, placeholdersToReplace *map[string]string) error {
	placeHolderValue, err := AskData(prompt.Prompt)
	if err != nil {
		return err
	}

	// Save asked value on the map
	(*placeholdersToReplace)[prompt.Placeholder] = placeHolderValue
	return nil
}

func numberValidationReplacePrompt(prompt Prompt, placeholdersToReplace *map[string]string) error {
	placeHolderValue, err := AskDataWithValidation(prompt.Prompt, validators.Number)
	if err != nil {
		return err
	}

	// Save asked value on the map
	(*placeholdersToReplace)[prompt.Placeholder] = placeHolderValue
	return nil
}

func numberRangeValidationReplacePrompt(prompt Prompt, placeholdersToReplace *map[string]string) error {
	validationParams := extractPromptValidationParamsValues(prompt.Validate)

	if len(validationParams) != 2 {
		return fmt.Errorf("the prompt structure of the %s placeholder has an incorrect validation parameters number", prompt.Placeholder)
	}

	// Get real parameters value type
	minValue := getRealPromptValidationParamType(validationParams[0])
	maxValue := getRealPromptValidationParamType(validationParams[1])

	// Check that all parametersValues are integers
	if reflect.TypeOf(minValue).Kind() != reflect.Int || reflect.TypeOf(maxValue).Kind() != reflect.Int {
		return fmt.Errorf("the prompt structure of the %s placeholder has an incorrect validation parameters type expected %s: %s received for minValue | expected %s: %s received for maxValue", prompt.Placeholder, reflect.TypeOf(minValue).Kind(), reflect.Int, reflect.TypeOf(maxValue).Kind(), reflect.Int)
	}

	// Check min max valid range
	if minValue.(int) > maxValue.(int) {
		return fmt.Errorf("the prompt structure of the %s placeholder has an incorrect validation parameters value, minValue could not be higher than maxValue", prompt.Placeholder)
	}

	validatorOptions := validators.NumberRangeOptions{
		Min: int32(minValue.(int)),
		Max: int32(maxValue.(int)),
	}

	placeHolderValue, err := AskDataWithValidation(prompt.Prompt, validators.NumberRange(validatorOptions))
	if err != nil {
		return err
	}

	// Save asked value on the map
	(*placeholdersToReplace)[prompt.Placeholder] = placeHolderValue
	return nil
}

func alphabetValidationReplacePrompt(prompt Prompt, placeholdersToReplace *map[string]string) error {
	validationParams := extractPromptValidationParamsValues(prompt.Validate)

	if len(validationParams) != 3 {
		return fmt.Errorf("the prompt structure of the %s placeholder has an incorrect validation parameters number", prompt.Placeholder)
	}

	// Get real parameters value type
	minValue := getRealPromptValidationParamType(validationParams[1])
	maxValue := getRealPromptValidationParamType(validationParams[2])
	caseSensitiveValue := getRealPromptValidationParamType(validationParams[0])

	// Check that all parametersValues are integers
	if reflect.TypeOf(caseSensitiveValue).Kind() != reflect.Int || reflect.TypeOf(minValue).Kind() != reflect.Int || reflect.TypeOf(maxValue).Kind() != reflect.Int {
		return fmt.Errorf("the prompt structure of the %s placeholder has an incorrect validation parameters type expected %s: %s received for caseSensitive | expected %s: %s received for minValue | expected %s: %s received for maxValue", prompt.Placeholder, reflect.TypeOf(caseSensitiveValue).Kind(), reflect.Int, reflect.TypeOf(minValue).Kind(), reflect.Int, reflect.TypeOf(maxValue).Kind(), reflect.Int)
	}

	// Check caseSensitive valid range
	if caseSensitiveValue.(int) < 0 || caseSensitiveValue.(int) > 2 {
		return fmt.Errorf("the prompt structure of the %s placeholder has an incorrect validation parameters value, caseSensitive only could be between 0 to 2", prompt.Placeholder)
	}

	// Check min max valid range
	if minValue.(int) <= maxValue.(int) {
		return fmt.Errorf("the prompt structure of the %s placeholder has an incorrect validation parameters value, minValue could not be higher than maxValue", prompt.Placeholder)
	}

	caseSensitiveType := validators.ParseCaseSensitiveIntValue(caseSensitiveValue.(int))
	stringLenghtOpts := validators.StringLenghtOptions{
		Min: int32(minValue.(int)),
		Max: int32(minValue.(int)),
	}

	placeHolderValue, err := AskDataWithValidation(prompt.Prompt, validators.Alphabet(caseSensitiveType, stringLenghtOpts))
	if err != nil {
		return err
	}

	// Save asked value on the map
	(*placeholdersToReplace)[prompt.Placeholder] = placeHolderValue
	return nil
}

func parseReplacePrompt(prompt Prompt, placeholdersToReplace *map[string]string) error {
	// Skip appname placeholder if exists, this placeholder is managed by MKS
	if prompt.Placeholder == config.PLACEHOLDER_APP_NAME {
		return nil
	}

	// Prompt -> number validation
	if strings.HasPrefix(prompt.Validate, config.VALIDATION_NUMBER) && !strings.HasPrefix(prompt.Validate, config.VALIDATION_NUMBER_RANGE) {
		if err := numberValidationReplacePrompt(prompt, placeholdersToReplace); err != nil {
			return err
		}
		return nil
	}

	// Prompt -> number range validation
	if strings.HasPrefix(prompt.Validate, config.VALIDATION_NUMBER_RANGE) {
		if err := numberRangeValidationReplacePrompt(prompt, placeholdersToReplace); err != nil {
			return err
		}
		return nil
	}

	// Prompt -> alpahet validation [az] | [AZ] | [azAZ]
	if strings.HasPrefix(prompt.Validate, config.VALIDATION_ALPHABET) {
		if err := alphabetValidationReplacePrompt(prompt, placeholdersToReplace); err != nil {
			return err
		}
		return nil
	}

	// Prompt -> no validation
	if err := noValidationReplacePrompt(prompt, placeholdersToReplace); err != nil {
		return err
	}
	return nil
}

func ParsePromptFile(promptFilePath string, placeholdersToReplace *map[string]string, extendsFilesToSkip *[]string) error {
	if global.Verbose {
		fmt.Printf("[+] Parsing %s prompt file..\n", promptFilePath)
	}

	templatePath := filepath.Dir(promptFilePath)

	// Read file content
	fileContent, err := ReadFile(promptFilePath)
	if err != nil {
		return err
	}

	// Variable to save parsed json data
	var parsedFile PromptsFileFormat

	// Parse json file and save data on parsedFile variable
	err = json.Unmarshal([]byte(fileContent), &parsedFile)
	if err != nil {
		return err
	}

	// if prompt file is empty
	if len(parsedFile.Prompts) == 0 {
		return nil
	}

	// Iterate every prompt
	for _, prompt := range parsedFile.Prompts {
		if prompt.Type == config.PROMPT_TYPE_REPLACE {
			parseReplacePrompt(prompt, placeholdersToReplace)
		}

		if prompt.Type == config.PROMPT_TYPE_EXTEND {
			extendFilePath := path.Join(templatePath, prompt.ExtendFile)

			// Validate if extendFile indicated on prompt struct is valid
			if !FileOrDirectoryExists(extendFilePath) {
				return fmt.Errorf("%s extend prompt has invalid extendFile: %s", filepath.Base(templatePath), prompt.ExtendFile)
			}

			// Ask user if wants to install extendFile
			installExtendFile, err := AskConfirm(prompt.Prompt)
			if err != nil {
				return err
			}

			// If user reply "no" append extendFile to skip list
			if !installExtendFile {
				*extendsFilesToSkip = append(*extendsFilesToSkip, prompt.ExtendFile)
			}
		}
	}

	return nil
}
