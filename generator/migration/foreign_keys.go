package migration

import (
	"fmt"

	"github.com/gobuffalo/fizz"
)

func parseForeignKeyRef(refs interface{}) (fizz.ForeignKeyRef, error) {
	fkr := fizz.ForeignKeyRef{}
	refMap, ok := refs.(map[string]interface{})
	if !ok {
		return fkr, fmt.Errorf(`invalid references format %s\nmust be "{"table": ["colum1", "column2"]}"`, refs)
	}
	if len(refMap) != 1 {
		return fkr, fmt.Errorf("only one table is supported as Foreign key reference")
	}
	for table, columns := range refMap {
		fkr.Table = table
		for _, c := range columns.([]interface{}) {
			fkr.Columns = append(fkr.Columns, fmt.Sprintf("%s", c))
		}
	}

	return fkr, nil
}
