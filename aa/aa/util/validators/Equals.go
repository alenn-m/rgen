package validators

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func PasswordEquals(input string) validation.RuleFunc {
	return func(value interface{}) error {
		if input != value.(string) {
			return errors.New("passwords are not the same")
		}

		return nil
	}
}
