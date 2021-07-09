package migration

import (
	"bytes"
	"fmt"
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

type parsedData struct {
	Name          string
	Fields        []parser.Field
	Relationships parser.Relationships
	Package       string
	Content       string
	Template      *template.Template
}

type Migration struct {
	Input *parsedData
}

func (m *Migration) Generate(input *parser.Parser, conf *config.Config) error {
	m.parseData(input, conf)

	attributes := []attrs.Attr{}
	for _, item := range m.Input.Fields {
		a, err := attrs.Parse(fmt.Sprintf("%s:%s", item.Key, item.Value))
		if err != nil {
			return err
		}

		attributes = append(attributes, a)
	}

	err := m.CreateMigration(&ctable.Options{
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

func (m *Migration) Save() error {
	err := files.MakeDirIfNotExist(dir)
	if err != nil {
		return err
	}

	migrationName := fmt.Sprintf("create_%s_table", m.Input.Name)
	goose.SetSequential(true)
	err = goose.CreateWithTemplate(nil, dir, m.Input.Template, migrationName, "sql")

	return err
}

func (m *Migration) GetContent() string {
	return m.Input.Content
}

func (m *Migration) parseData(input *parser.Parser, conf *config.Config) {
	m.Input = &parsedData{
		Name:          input.Name,
		Fields:        input.Fields,
		Relationships: input.Relationships,
		Package:       conf.Package,
	}
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

	for modelItem, relationhip := range m.Input.Relationships {
		if relationhip == model.BelongsTo {
			fk := fmt.Sprintf("%sID", modelItem)
			err = t.ForeignKey(fk, map[string]interface{}{
				inflection.Plural(modelItem): []interface{}{fk},
			}, nil)
			if err != nil {
				return err
			}
		}
	}

	err = m.generateGooseMigration(opts, t)
	if err != nil {
		return err
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

	var tpl bytes.Buffer
	err = sqlMigrationTemplate.Execute(&tpl, nil)
	if err != nil {
		return err
	}

	m.Input.Content = tpl.String()
	m.Input.Template = sqlMigrationTemplate

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
