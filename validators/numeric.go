package validators

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

const int32minlengh = -2147483648
const int32maxlengh = 2147483647

type NumberRangeOptions struct {
	Min int32
	Max int32
}

// ValidateIntInRange is a custom validation function to ensure that the input is an integer within the specified range.
func NumberRange(opt ...NumberRangeOptions) survey.Validator {
	return func(val interface{}) error {
		strVal, ok := val.(string)
		if !ok {
			return errors.New("the value is not a string")
		}

		intVal, err := strconv.Atoi(strVal)
		if err != nil {
			return errors.New("the value is not a valid number")
		}

		var options NumberRangeOptions

		options = NumberRangeOptions{
			Min: int32minlengh,
			Max: int32maxlengh,
		}

		// Check if options are custom
		if len(opt) > 0 {
			options = opt[0]
		}

		if intVal < int(options.Min) || intVal > int(options.Max) {
			return fmt.Errorf("the number must be in the range of %d to %d", options.Min, options.Max)
		}

		return nil
	}
}

func Number(val interface{}) error {
	strVal, ok := val.(string)
	if !ok {
		return errors.New("the value is not a string")
	}

	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		return errors.New("the value is not a valid number")
	}

	// Check if intVal overflows Int32
	if intVal < int32minlengh || intVal > int32maxlengh {
		return fmt.Errorf("the number cannot be lower than %d or higher than %d", int32minlengh, int32maxlengh)
	}

	return nil
}
