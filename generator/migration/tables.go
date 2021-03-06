package migration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gobuffalo/fizz"
)

// Table is the table definition for fizz.
type Table struct {
	Name              string `db:"name"`
	Columns           []fizz.Column
	Indexes           []fizz.Index
	ForeignKeys       []fizz.ForeignKey
	primaryKeys       []string
	Options           map[string]interface{}
	columnsCache      map[string]struct{}
	useTimestampMacro bool
}

func (t Table) String() string {
	return t.Fizz()
}

// Fizz returns the fizz DDL to create the table.
func (t Table) Fizz() string {
	var buff bytes.Buffer
	timestampsOpt, _ := t.Options["timestamps"].(bool)
	// Write table options
	o := make([]string, 0, len(t.Options))
	for k, v := range t.Options {
		vv, _ := json.Marshal(v)
		o = append(o, fmt.Sprintf("%s: %s", k, string(vv)))
	}
	if len(o) > 0 {
		sort.SliceStable(o, func(i, j int) bool { return o[i] < o[j] })
		buff.WriteString(fmt.Sprintf("create_table(\"%s\", {%s}) {\n", t.Name, strings.Join(o, ", ")))
	} else {
		buff.WriteString(fmt.Sprintf("create_table(\"%s\") {\n", t.Name))
	}
	// Write columns
	if t.useTimestampMacro {
		for _, c := range t.Columns {
			if c.Name == "created_at" || c.Name == "updated_at" {
				continue
			}
			buff.WriteString(fmt.Sprintf("\t%s\n", c.String()))
		}
	} else {
		for _, c := range t.Columns {
			buff.WriteString(fmt.Sprintf("\t%s\n", c.String()))
		}
	}
	if t.useTimestampMacro {
		buff.WriteString("\tt.Timestamps()\n")
	} else if timestampsOpt {
		// Missing timestamp columns will only be added on fizz execution, so we need to consider them as present.
		if !t.HasColumns("created_at") {
			buff.WriteString(fmt.Sprintf("\t%s\n", fizz.CREATED_COL.String()))
		}
		if !t.HasColumns("updated_at") {
			buff.WriteString(fmt.Sprintf("\t%s\n", fizz.UPDATED_COL.String()))
		}
	}
	// Write primary key (single column pk will be written in inline form as the column opt)
	if len(t.primaryKeys) > 1 {
		pks := make([]string, len(t.primaryKeys))
		for i, pk := range t.primaryKeys {
			pks[i] = fmt.Sprintf("\"%s\"", pk)
		}
		buff.WriteString(fmt.Sprintf("\tt.PrimaryKey(%s)\n", strings.Join(pks, ", ")))
	}
	// Write indexes
	for _, i := range t.Indexes {
		buff.WriteString(fmt.Sprintf("\t%s\n", i.String()))
	}
	// Write foreign keys
	for _, fk := range t.ForeignKeys {
		buff.WriteString(fmt.Sprintf("\t%s\n", fk.String()))
	}
	buff.WriteString("}")
	return buff.String()
}

// HasColumns checks if the Table has all the given columns.
func (t *Table) HasColumns(args ...string) bool {
	for _, a := range args {
		if _, ok := t.columnsCache[a]; !ok {
			return false
		}
	}
	return true
}

// NewTable creates a new Table.
func NewTable(name string, opts map[string]interface{}) Table {
	if opts == nil {
		opts = make(map[string]interface{})
	}

	return Table{
		Name:         name,
		Columns:      []fizz.Column{},
		Indexes:      []fizz.Index{},
		Options:      opts,
		columnsCache: map[string]struct{}{},
	}
}

// Column adds a column to the table definition.
func (t *Table) Column(name string, colType string, options fizz.Options) error {
	if _, found := t.columnsCache[name]; found {
		return fmt.Errorf("duplicated column %s", name)
	}
	var primary bool
	if _, ok := options["primary"]; ok {
		if t.primaryKeys != nil {
			return errors.New("could not define multiple primary keys")
		}
		primary = true
		t.primaryKeys = []string{name}
	}
	c := fizz.Column{
		Name:    name,
		ColType: colType,
		Options: options,
		Primary: primary,
	}
	if t.columnsCache == nil {
		t.columnsCache = make(map[string]struct{})
	}
	t.columnsCache[name] = struct{}{}
	// Ensure id is first
	if name == "id" {
		t.Columns = append([]fizz.Column{c}, t.Columns...)
	} else {
		t.Columns = append(t.Columns, c)
	}
	if (name == "created_at" || name == "updated_at") && colType != "timestamp" {
		// timestamp macro only works for time type
		t.useTimestampMacro = false
	}
	return nil
}

// UnFizz returns the fizz DDL to remove the table.
func (t Table) UnFizz() string {
	return fmt.Sprintf("drop_table(\"%s\")", t.Name)
}

// Timestamp is a shortcut to add a timestamp column with default options.
func (t *Table) Timestamp(name string) error {
	return t.Column(name, "timestamp", fizz.Options{})
}

// ForeignKey adds a new foreign key to the table definition.
func (t *Table) ForeignKey(column string, refs interface{}, options fizz.Options) error {
	fkr, err := parseForeignKeyRef(refs)
	if err != nil {
		return err
	}
	fk := fizz.ForeignKey{
		Column:     column,
		References: fkr,
		Options:    options,
	}

	if options["name"] != nil {
		var ok bool
		fk.Name, ok = options["name"].(string)
		if !ok {
			return fmt.Errorf(`expected options field "name" to be of type "string" but got "%T"`, options["name"])
		}
	} else {
		fk.Name = fmt.Sprintf("%s_%s_%s_fk", t.Name, fk.References.Table, strings.Join(fk.References.Columns, "_"))
	}

	t.ForeignKeys = append(t.ForeignKeys, fk)
	return nil
}
