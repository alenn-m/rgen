package migration

import (
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/alenn-m/rgen/generator/model"
	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/files"
	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/gobuffalo/pop/v5/genny/fizz/ctable"
	"github.com/jinzhu/inflection"
	"github.com/pressly/goose"
)

var dir = "database/migrations"

type Input struct {
	Name          string
	Fields        []parser.Field
	Relationships parser.Relationships
}

type Migration struct {
	Input  *Input
	Config *config.Config
}

func (m *Migration) Init(input *Input, conf *config.Config) {
	m.Input = input
	m.Config = conf
}

func (m *Migration) Generate() error {
	err := files.MakeDirIfNotExist(dir)
	if err != nil {
		return err
	}

	attributes := []attrs.Attr{}
	for _, item := range m.Input.Fields {
		a, err := attrs.Parse(fmt.Sprintf("%s:%s", item.Key, item.Value))
		if err != nil {
			return err
		}

		attributes = append(attributes, a)
	}

	err = m.CreateMigration(&ctable.Options{
		TableName:  m.Input.Name,
		Path:       dir,
		Type:       "sql",
		Attrs:      attributes,
		Translator: translators.NewMySQL("", ""),
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *Migration) CreateMigration(opts *ctable.Options) error {
	t := NewTable(inflection.Plural(opts.TableName), map[string]interface{}{
		"timestamps": false,
	})

	err := t.Column(fmt.Sprintf("%sID", opts.TableName), "integer", fizz.Options{"primary": true})
	if err != nil {
		return err
	}

	for _, attr := range opts.Attrs {
		colType := m.fizzColType(attr.CommonType())

		if err := t.Column(attr.Name.String(), colType, fizz.Options{}); err != nil {
			return err
		}
	}

	err = t.Column("CreatedAt", "timestamp", fizz.Options{"default_raw": "CURRENT_TIMESTAMP"})
	if err != nil {
		return err
	}

	err = t.Column("UpdatedAt", "timestamp", fizz.Options{"default_raw": "CURRENT_TIMESTAMP"})
	if err != nil {
		return err
	}

	m2mTables := []Table{}

	for modelItem, relationhip := range m.Input.Relationships {
		switch relationhip {
		case model.BelongsTo:
			fk := fmt.Sprintf("%sID", modelItem)
			err = t.ForeignKey(fk, map[string]interface{}{
				inflection.Plural(modelItem): []interface{}{fk},
			}, nil)
			if err != nil {
				return err
			}
		case model.ManyToMany:
			// create slice of joining tables
			tables := []string{m.Input.Name, modelItem}
			// sort them
			sort.Strings(tables)
			// create the joining table name
			r := inflection.Plural(fmt.Sprintf("%s%s", tables[0], tables[1]))

			m2m := NewTable(r, map[string]interface{}{
				"timestamps": false,
			})

			fk1 := fmt.Sprintf("%sID", tables[0])
			fk2 := fmt.Sprintf("%sID", tables[1])

			m2m.Column(fk1, "integer", fizz.Options{"primary": true})
			m2m.ForeignKey(fk1, map[string]interface{}{
				inflection.Plural(tables[0]): []interface{}{fk1},
			}, nil)
			m2m.Column(fk2, "integer", fizz.Options{"primary": true})
			m2m.ForeignKey(fk2, map[string]interface{}{
				inflection.Plural(tables[1]): []interface{}{fk2},
			}, nil)

			m2mTables = append(m2mTables, m2m)
		default:
			return fmt.Errorf("invalid relationship: %s", relationhip)
		}
	}

	err = m.generateGooseMigration(opts, t)
	if err != nil {
		return err
	}

	for _, m2m := range m2mTables {
		err = m.generateGooseMigration(opts, m2m)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migration) generateGooseMigration(opts *ctable.Options, t Table) error {
	upString, err := fizz.AString(t.Fizz(), opts.Translator)
	if err != nil {
		return err
	}

	downString, err := fizz.AString(t.UnFizz(), opts.Translator)
	if err != nil {
		return err
	}

	sqlMigrationTemplate := template.Must(template.New("goose.sql-migration").Parse(fmt.Sprintf(`-- +goose Up
%s

-- +goose Down
%s
`, upString, downString)))

	migrationName := fmt.Sprintf("create_%s_table", t.Name)
	goose.SetSequential(true)
	err = goose.CreateWithTemplate(nil, dir, sqlMigrationTemplate, migrationName, "sql")
	if err != nil {
		return err
	}

	return nil
}

func (m *Migration) fizzColType(s string) string {
	switch strings.ToLower(s) {
	case "int", "int8", "int16", "int32", "int64":
		return "integer"
	case "time", "datetime":
		return "timestamp"
	case "uuid.uuid", "uuid":
		return "uuid"
	case "nulls.float32", "nulls.float64":
		return "float"
	case "slices.string", "slices.uuid", "[]string":
		return "varchar[]"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "numeric[]"
	case "slices.int":
		return "int[]"
	case "slices.map":
		return "jsonb"
	case "float", "float32", "float64":
		return "decimal"
	case "blob", "[]byte":
		return "blob"
	default:
		if strings.HasPrefix(s, "nulls.") {
			return m.fizzColType(strings.Replace(s, "nulls.", "", -1))
		}
		return strings.ToLower(s)
	}
}
