package model

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

var dir = "models"

type Input struct {
	Name   string
	Fields []parser.Field
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
	m.parseFields()

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

func (m *Model) parseFields() {
	fields := fmt.Sprintf("ID int64 `json:\"id\" orm:\"pk\"`\n")
	for _, item := range m.Input.Fields {
		fields += fmt.Sprintf("%s %s `json:\"%s\"`\n", strcase.ToCamel(item.Key), item.Value, strcase.ToSnake(item.Key))
	}

	m.ParsedData.Fields = fields
}

func (m *Model) saveFile(content []byte, location string) error {
	err := ioutil.WriteFile(location, content, 0644)

	return err
}
