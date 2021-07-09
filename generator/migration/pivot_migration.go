package migration

import (
	"fmt"
	"sort"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/gobuffalo/pop/v5/genny/fizz/ctable"
	"github.com/jinzhu/inflection"
)

type PivotMigrationEntry struct {
	TableOne string
	TableTwo string
}

type PivotMigration struct {
	Tables []PivotMigrationEntry
}

func (p *PivotMigration) Init(tables []PivotMigrationEntry) {
	p.Tables = tables
}

func (p *PivotMigration) Generate() error {
	m2mTables := map[string]Table{}

	for _, item := range p.Tables {
		// create slice of joining tables
		tables := []string{item.TableOne, item.TableTwo}
		// sort them
		sort.Strings(tables)
		// create the joining table name
		r := inflection.Plural(fmt.Sprintf("%s%s", tables[0], tables[1]))

		m2m := NewTable(r, map[string]interface{}{
			"timestamps": false,
		})

		fk1 := fmt.Sprintf("%sID", tables[0])
		fk2 := fmt.Sprintf("%sID", tables[1])

		m2m.Column(fk1, "integer", fizz.Options{})
		m2m.ForeignKey(fk1, map[string]interface{}{
			inflection.Plural(tables[0]): []interface{}{fk1},
		}, nil)
		m2m.Column(fk2, "integer", fizz.Options{})
		m2m.ForeignKey(fk2, map[string]interface{}{
			inflection.Plural(tables[1]): []interface{}{fk2},
		}, nil)
		m2m.PrimaryKey(fk1, fk2)

		m2mTables[r] = m2m
	}

	opts := &ctable.Options{
		Path:       dir,
		Type:       "sql",
		Translator: translators.NewMySQL("", ""),
	}
	for table, m2m := range m2mTables {
		m := &Migration{Input: &parsedData{
			Name: table,
		}}

		err := m.generateGooseMigration(opts, m2m)
		if err != nil {
			return err
		}

		err = m.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
