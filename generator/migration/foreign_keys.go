package migration

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gobuffalo/fizz"
)

// ForeignKey
type ForeignKey struct {
	Name       string
	Column     string
	References fizz.ForeignKeyRef
	Options    fizz.Options
}

func (f ForeignKey) String() string {
	refs := fmt.Sprintf(`{"%s": ["%s"]}`, f.References.Table, strings.Join(f.References.Columns, `", "`))
	var opts map[string]interface{}
	if f.Options == nil {
		opts = make(map[string]interface{}, 0)
	} else {
		opts = f.Options
	}

	o := make([]string, 0, len(opts))
	for k, v := range opts {
		vv, _ := json.Marshal(v)
		o = append(o, fmt.Sprintf("%s: %s", k, string(vv)))
	}
	sort.SliceStable(o, func(i, j int) bool { return o[i] < o[j] })
	return fmt.Sprintf(`t.ForeignKey("%s", %s, {%s})`, f.Column, refs, strings.Join(o, ", "))
}

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
