package model

import (
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

const BelongsTo = "belongsTo"
const HasMany = "hasMany"
const ManyToMany = "manyToMany"

var dir = "models"

type Input struct {
	Name          string
	Relationships map[string]string
	Fields        []parser.Field
}

type Model struct {
	Input  *Input
	Config *config.Config

	ParsedData parsedData
}

type parsedData struct {
	Name   string
	Fields string
}

func (m *Model) Init(input *Input, conf *config.Config) {
	m.Input = input
	m.Config = conf
}

func (m *Model) Generate() error {
	m.parseModelName()
	err := m.parseFields()
	if err != nil {
		return err
	}

	contentString := TEMPLATE
	contentString = strings.Replace(contentString, "{{Model}}", m.ParsedData.Name, -1)
	contentString = strings.Replace(contentString, "{{Fields}}", m.ParsedData.Fields, -1)
	contentString = strings.Replace(contentString, "{{Root}}", m.Config.Package, -1)

	location, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	content, err := format.Source([]byte(contentString))
	if err != nil {
		return err
	}

	err = m.saveFile(content, fmt.Sprintf("%s/%s.go", location, m.Input.Name))
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) parseModelName() {
	m.ParsedData.Name = inflection.Singular(strings.Title(strings.ToLower(m.Input.Name)))
}

func (m *Model) parseFields() error {
	fields := fmt.Sprintf("ID int64 `json:\"id\" orm:\"pk\"`\n")
	for _, item := range m.Input.Fields {
		fields += fmt.Sprintf("%s %s `json:\"%s\"`\n", strcase.ToCamel(item.Key), item.Value, strcase.ToSnake(item.Key))
	}

	for key, relationship := range m.Input.Relationships {
		switch relationship {
		case BelongsTo:
			fields += fmt.Sprintf("%s *%s `json:\"%s\"`\n", key, key, strcase.ToSnake(key))
		case HasMany:
			fields += fmt.Sprintf("%s []%s `json:\"%s\"`\n", inflection.Plural(key), key, strcase.ToSnake(inflection.Plural(key)))
		case ManyToMany:
			// create slice of joining tables
			tables := []string{strings.ToLower(m.Input.Name), strings.ToLower(key)}
			// sort them
			sort.Strings(tables)
			// create the joining table name
			r := inflection.Plural(fmt.Sprintf("%s_%s", tables[0], tables[1]))
			fields += fmt.Sprintf("%s []%s `json:\"%s\" gorm:\"many2many:%s\"`\n", inflection.Plural(key), key, strcase.ToSnake(inflection.Plural(key)), r)
		default:
			return errors.New("invalid relationship")
		}
	}

	m.ParsedData.Fields = fields

	return nil
}

func (m *Model) saveFile(content []byte, location string) error {
	err := ioutil.WriteFile(location, content, 0644)

	return err
}
