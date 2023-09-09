package validators

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
)

// Enumeración para los valores permitidos
type CaseType int

const (
	Mayus CaseType = iota
	Minus
	None
)

type StringLenghtOptions struct {
	Min int32
	Max int32
}

const minStringLenght = 0
const maxStringLenght = 65536

func ParseCaseSensitiveIntValue(value int) CaseType {
	switch value {
	case 0:
		return Mayus
	case 1:
		return Minus
	default:
		return None
	}
}

func Alphabet(caseSensitive CaseType, opt ...StringLenghtOptions) survey.Validator {
	return func(val interface{}) error {
		str, ok := val.(string)
		if !ok {
			return errors.New("the provided value is not a string")
		}

		// Regular expression to verify that the string contains only A to Z letters (uppercase or lowercase)
		regex := "^[a-zA-Z]+$"

		if caseSensitive == Mayus {
			regex = "^[A-Z]+$"
		}

		if caseSensitive == Minus {
			regex = "^[a-z]+$"
		}

		options := StringLenghtOptions{
			Min: minStringLenght,
			Max: maxStringLenght,
		}

		// Check if options are custom
		if len(opt) > 0 {
			if opt[0].Min < minStringLenght {
				opt[0].Min = minStringLenght
			}

			if opt[0].Max > maxStringLenght {
				opt[0].Max = maxStringLenght
			}

			options = opt[0]
		}

		// Validate string lenght range
		if len(str) < int(options.Min) {
			return fmt.Errorf("the string lenght cannot be lower than %d", options.Min)
		}

		if len(str) > int(options.Max) {
			return fmt.Errorf("the string lenght cannot be higher than %d", options.Max)
		}

		match, err := regexp.MatchString(regex, str)
		if err != nil {
			return err
		}

		if !match {
			return errors.New("the string contains invalid characters")
		}

		return nil
	}
}
