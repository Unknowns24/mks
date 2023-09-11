package utils

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func AskConfirm(question string) (bool, error) {
	var answer bool

	prompt := &survey.Confirm{
		Message: question,
		Default: false, // Optional: Set default value to false
	}

	err := survey.AskOne(prompt, &answer)

	if err != nil {
		return false, err
	}

	return answer, nil
}

func AskData(question string) (string, error) {
	response := ""
	confirmed := false

	prompt := &survey.Input{
		Message: question,
	}

	err := survey.AskOne(prompt, &response)

	if err != nil {
		return "", err
	}

	// Confirmate value
	confirm := &survey.Confirm{
		Message: fmt.Sprintf("Your answer was \"%s\", Do you want to continue?", response),
		Default: false, // Optional: Set default value to false
	}

	// Ask if value is correct
	err = survey.AskOne(confirm, &confirmed)
	if err != nil {
		return "", err
	}

	if confirmed {
		return response, nil
	} else {
		return AskData(question)
	}
}

func AskDataWithValidation(questionTitle string, validator survey.Validator) (string, error) {
	response := ""
	confirmed := false

	prompt := &survey.Input{
		Message: questionTitle,
	}

	err := survey.AskOne(prompt, &response, survey.WithValidator(validator))

	if err != nil {
		return "", err
	}

	// Confirmate value
	confirm := &survey.Confirm{
		Message: fmt.Sprintf("Your answer was \"%s\", Do you want to continue?", response),
		Default: false, // Optional: Set default value to false
	}
	// Ask if value is correct
	err = survey.AskOne(confirm, &confirmed)
	if err != nil {
		return "", err
	}

	if confirmed {
		return response, nil
	} else {
		return AskDataWithValidation(questionTitle, validator)
	}
}
