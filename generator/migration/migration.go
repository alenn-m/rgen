package migration

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/alenn-m/rgen/v2/generator/model"
	"github.com/alenn-m/rgen/v2/generator/parser"
	"github.com/alenn-m/rgen/v2/util/config"
	"github.com/alenn-m/rgen/v2/util/files"
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
	Sequential    bool
	Content       string
	Template      *template.Template
}

// Migration generator
type Migration struct {
	parsedData *parsedData
}

// Generate generates migration
func (m *Migration) Generate(input *parser.Parser, conf *config.Config) error {
	m.parseData(input, conf)

	attributes := []attrs.Attr{}
	for _, item := range m.parsedData.Fields {
		a, err := attrs.Parse(fmt.Sprintf("%s:%s", item.Key, item.Value))
		if err != nil {
			return err
		}

		attributes = append(attributes, a)
	}

	err := m.createMigration(&ctable.Options{
		TableName:  m.parsedData.Name,
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

// Save saves generated migration to file
func (m *Migration) Save() error {
	err := files.MakeDirIfNotExist(dir)
	if err != nil {
		return err
	}

	migrationName := fmt.Sprintf("create_%s_table", m.parsedData.Name)
	goose.SetSequential(m.parsedData.Sequential)
	err = goose.CreateWithTemplate(nil, dir, m.parsedData.Template, migrationName, "sql")

	return err
}

// GetContent return generated migration content
func (m *Migration) GetContent() string {
	return m.parsedData.Content
}

func (m *Migration) parseData(input *parser.Parser, conf *config.Config) {
	m.parsedData = &parsedData{
		Name:          input.Name,
		Fields:        input.Fields,
		Relationships: input.Relationships,
		Package:       conf.Package,
		Sequential:    conf.Migration.Sequential,
	}
}

func (m *Migration) createMigration(opts *ctable.Options) error {
	t := NewTable(inflection.Plural(opts.TableName), map[string]interface{}{
		"timestamps": false,
	})

	err := t.Column(fmt.Sprintf("%sID", opts.TableName), "integer", fizz.Options{"primary": true})
	if err != nil {
		return err
	}

	for _, attr := range opts.Attrs {
		colType := strings.ToLower(attr.CommonType())

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

	for modelItem, relationhip := range m.parsedData.Relationships {
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

	m.parsedData.Content = tpl.String()
	m.parsedData.Template = sqlMigrationTemplate

	return nil
}
