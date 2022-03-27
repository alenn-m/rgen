package validators

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/inflection"
)

func Unique(table, field string) validation.RuleFunc {
	return func(value interface{}) error {
		exists := false

		_ = db.Get(&exists, fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE %s = ?)", table, field), value)
		if exists {
			return fmt.Errorf("%s in table %s is already taken", value, table)
		}

		return nil
	}
}

func UniqueExcept(table, field string, id int64) validation.RuleFunc {
	tableSingular := inflection.Singular(table)
	idColumn := fmt.Sprintf("%sID", tableSingular)

	return func(value interface{}) error {
		exists := false

		_ = db.Get(&exists, fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE %s = ? AND %s != ?)", table, field, idColumn), value, id)
		if exists {
			return fmt.Errorf("%s in table %s is already taken", value, table)
		}

		return nil
	}
}
