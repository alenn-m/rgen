package migration

import (
	"fmt"
	"strings"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/gobuffalo/pop/v5/genny/fizz/ctable"
)

var dir = "database/migrations"

type Input struct {
	Name   string
	Fields []parser.Field
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
	attributes := []attrs.Attr{}
	for _, item := range m.Input.Fields {
		a, err := attrs.Parse(fmt.Sprintf("%s:%s", item.Key, item.Value))
		if err != nil {
			return err
		}

		attributes = append(attributes, a)
	}

	up, down, err := m.CreateMigration(&ctable.Options{
		TableName:  m.Input.Name,
		Path:       dir,
		Type:       "sql",
		Attrs:      attributes,
		Translator: translators.NewMySQL("", ""),
	})
	if err != nil {
		return err
	}

	fmt.Println(up)
	fmt.Println(down)

	return nil
}

func (m *Migration) CreateMigration(opts *ctable.Options) (string, string, error) {
	if err := opts.Validate(); err != nil {
		return "", "", err
	}

	t := NewTable(opts.TableName, map[string]interface{}{
		"timestamps": false,
	})

	for _, attr := range opts.Attrs {
		o := fizz.Options{}
		name := attr.Name.String()
		colType := m.fizzColType(attr.CommonType())
		if name == "id" {
			o["primary"] = true
		}
		if strings.HasPrefix(attr.GoType(), "nulls.") {
			o["null"] = true
		}
		if err := t.Column(name, colType, o); err != nil {
			return "", "", err
		}
	}

	err := t.Timestamp("CreatedAt")
	if err != nil {
		return "", "", err
	}
	err = t.Timestamp("UpdatedAt")
	if err != nil {
		return "", "", err
	}

	up := t.Fizz()
	down := t.UnFizz()

	upString, err := fizz.AString(up, opts.Translator)
	if err != nil {
		return "", "", err
	}

	downString, err := fizz.AString(down, opts.Translator)
	if err != nil {
		return "", "", err
	}

	return upString, downString, nil
}

func (m *Migration) fizzColType(s string) string {
	switch strings.ToLower(s) {
	case "int":
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
	case "float32", "float64", "float":
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
