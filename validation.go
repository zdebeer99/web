package webapp

import (
	"fmt"
)

//
// This Section contains basic validation functions as the bulk of validation
// should be done by a third party app like https://github.com/asaskevich/govalidator
//

func ValidateString(value string, minlength, maxlength int) error {
	size := len(value)
	if size < minlength {
		return fmt.Errorf("Value must contain equal or more than %v characters.", minlength)
	}
	if size > maxlength {
		return fmt.Errorf("Value must contain less than %v characters.", maxlength)
	}
	return nil
}

func ValidateInt(value, min, max int) error {
	if value < min {
		return fmt.Errorf("Value must contain equal or more than %v characters.", min)
	}
	if value > max {
		return fmt.Errorf("Value must contain less than %v characters.", max)
	}
	return nil
}
