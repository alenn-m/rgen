package model

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/templates"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

//go:embed "template.tmpl"
var TEMPLATE string

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

	ParsedModelData     parsedModelData
	ParsedMigrationData parsedMigrationData
}

type parsedMigrationData struct {
	Models string
	Root   string
}

type parsedModelData struct {
	Model  string
	Fields string
}

func (m *Model) Init(input *Input, conf *config.Config) {
	m.Input = input
	m.Config = conf
}

func (m *Model) Generate() error {
	err := m.parseData()
	if err != nil {
		return err
	}

	p, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	content, err := templates.ParseTemplate(TEMPLATE, m.ParsedModelData, nil)
	if err != nil {
		return err
	}

	err = m.createFile(p, content)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) parseData() error {
	fields := fmt.Sprintf("ID %sID `json:\"id\" db:\"%s\"`\n", m.Input.Name, strcase.ToCamel(m.Input.Name))
	for _, item := range m.Input.Fields {
		camelName := strcase.ToCamel(item.Key)
		snakeName := strcase.ToSnake(item.Key)

		itemType := item.Value
		if strings.ToLower(item.Value) == "id" {
			itemType = fmt.Sprintf("%sID", camelName)
		}

		fields += fmt.Sprintf("%s %s `json:\"%s\"`\n", camelName, itemType, snakeName)
	}

	for key, relationship := range m.Input.Relationships {
		switch relationship {
		case BelongsTo:
			fields += fmt.Sprintf("%s *%s `json:\"%s\"`\n", key, key, strcase.ToSnake(key))
		case HasMany:
			fields += fmt.Sprintf("%s []%s `json:\"%s\"`\n", inflection.Plural(key), key, strcase.ToSnake(inflection.Plural(key)))
		case ManyToMany:
			fields += fmt.Sprintf("%s []%s `json:\"%s\"`\n", inflection.Plural(key), key, strcase.ToSnake(key))
		default:
			return errors.New("invalid relationship")
		}
	}

	m.ParsedModelData = parsedModelData{
		Model:  inflection.Singular(strings.Title(m.Input.Name)),
		Fields: fields,
	}

	return nil
}

func (m *Model) createFile(location, content string) error {
	filename := strings.Title(strings.ToLower(m.Input.Name))
	err := ioutil.WriteFile(fmt.Sprintf("%s/%s.go", location, filename), []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
