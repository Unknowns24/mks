package validators

import (
	"errors"
	"strconv"
)

// ValidateIntInRange is a custom validation function to ensure that the input is an integer within the specified range.
func ValidatePortRange(val interface{}) error {
	strVal, ok := val.(string)
	if !ok {
		return errors.New("the value is not a string")
	}

	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		return errors.New("the value is not a valid number")
	}

	if intVal < 1024 || intVal > 65535 {
		return errors.New("the number must be in the range of 1024 to 65535")
	}

	return nil
}
