package validators

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/inflection"
)

func RecordExist(table string) validation.RuleFunc {
	tableSingular := inflection.Singular(table)
	idColumn := fmt.Sprintf("%sID", tableSingular)

	return func(id interface{}) error {
		exists := false

		err := db.Get(&exists, fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE %s = ?)", table, idColumn), id)
		if err != nil || !exists {
			return fmt.Errorf("%s with ID = %d doesn't exist", tableSingular, id)
		}

		return nil
	}
}
