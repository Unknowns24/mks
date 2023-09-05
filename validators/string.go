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
			Min: 0,
			Max: 65536,
		}

		// Check if options are custom
		if len(opt) > 0 {
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
