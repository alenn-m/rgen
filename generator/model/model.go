package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/files"
	"github.com/alenn-m/rgen/util/templates"
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

	output, err := ioutil.ReadFile(fmt.Sprintf("%s/src/github.com/alenn-m/rgen/generator/model/template.tmpl", os.Getenv("GOPATH")))
	if err != nil {
		return err
	}

	content, err := templates.ParseTemplate(string(output), m.ParsedModelData, nil)
	if err != nil {
		return err
	}

	err = m.createFile(p, content, "model.go")
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) parseData() error {
	fields := fmt.Sprintf("ID %sID `json:\"id\" gorm:\"primaryKey\"`\n", m.Input.Name)
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

	m.ParsedModelData = parsedModelData{
		Model:  inflection.Singular(strings.Title(m.Input.Name)),
		Fields: fields,
	}

	return nil
}

func (m *Model) createFile(location, content, filename string) error {
	servicePath := fmt.Sprintf("%s/%s", location, strings.ToLower(m.Input.Name))

	err := files.MakeDirIfNotExist(servicePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", location, filename), []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
