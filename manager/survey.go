package manager

import "github.com/AlecAivazis/survey/v2"

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
		Message: response + ", Is it correct?",
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

func AskDataWithValidation(questionTitle, validationTitle string, validator survey.Validator) (string, error) {
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
		Message: response + ", Is it correct?",
	}

	// Ask if value is correct
	err = survey.AskOne(confirm, &confirmed)
	if err != nil {
		return "", err
	}

	if confirmed {
		return response, nil
	} else {
		return AskDataWithValidation(questionTitle, validationTitle, validator)
	}
}
