package validators

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/inflection"
	"github.com/jmoiron/sqlx"
)

func RecordsExist(table string) validation.RuleFunc {
	tableSingular := inflection.Singular(table)
	idColumn := fmt.Sprintf("%sID", tableSingular)

	return func(ids interface{}) error {
		exists := false

		q, args, err := sqlx.In(fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE %s IN (?))", table, idColumn), ids)
		if err != nil {
			return err
		}

		err = db.Get(&exists, q, args...)
		if err != nil || !exists {
			return errors.New("one of the items doesn't exists")
		}

		return nil
	}
}
