package model

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alenn-m/rgen/v2/generator/parser"
	"github.com/alenn-m/rgen/v2/util/config"
	"github.com/alenn-m/rgen/v2/util/files"
	"github.com/alenn-m/rgen/v2/util/templates"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

//go:embed "template.tmpl"
var TEMPLATE string

// ErrInvalidRelationship - returned when relationship is not valid
var ErrInvalidRelationship = errors.New("invalid relationship")

// BelongsTo - belongsToMany relationship
const BelongsTo = "belongsTo"

// HasMany - hasMany relationship
const HasMany = "hasMany"

// ManyToMany - manyToMany relationship
const ManyToMany = "manyToMany"

var dir = "models"

// Model generator
type Model struct {
	parsedData parsedData
}

// Save saves the generated content to file
func (m *Model) Save() error {
	return m.createFile(m.GetContent())
}

// GetContent content getter
func (m *Model) GetContent() string {
	return m.parsedData.Content
}

type parsedData struct {
	Model   string
	Fields  string
	Package string
	Content string
}

// Generate generates the '{MODEL}.go' file
func (m *Model) Generate(input *parser.Parser, conf *config.Config) error {
	err := m.parseData(input, conf)
	if err != nil {
		return err
	}

	content, err := templates.ParseTemplate(TEMPLATE, m.parsedData, nil)
	if err != nil {
		return err
	}

	m.parsedData.Content = content

	return nil
}

func (m *Model) parseData(input *parser.Parser, conf *config.Config) error {
	fields := fmt.Sprintf("ID %sID `json:\"id\" db:\"%sID\"`\n", input.Name, strcase.ToCamel(input.Name))
	for _, item := range input.Fields {
		camelName := strcase.ToCamel(item.Key)
		snakeName := strcase.ToSnake(item.Key)

		fields += fmt.Sprintf("%s %s `json:\"%s\"`\n", camelName, item.Value, snakeName)
	}

	relationshipKeys := []string{}
	for item := range input.Relationships {
		relationshipKeys = append(relationshipKeys, item)
	}

	sort.Strings(relationshipKeys)

	for _, key := range relationshipKeys {
		relationship := input.Relationships[key]
		switch relationship {
		case BelongsTo:
			fields += fmt.Sprintf("%s *%s `json:\"%s\"`\n", key, key, strcase.ToSnake(key))
		case HasMany:
			fields += fmt.Sprintf("%s []%s `json:\"%s\"`\n", inflection.Plural(key), key, strcase.ToSnake(inflection.Plural(key)))
		case ManyToMany:
			fields += fmt.Sprintf("%s []%s `json:\"%s\"`\n", inflection.Plural(key), key, strcase.ToSnake(key))
		default:
			return ErrInvalidRelationship
		}
	}

	m.parsedData = parsedData{
		Model:   inflection.Singular(strings.Title(input.Name)),
		Fields:  fields,
		Package: conf.Package,
	}

	return nil
}

func (m *Model) createFile(content string) error {
	location, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	err = files.MakeDirIfNotExist(location)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.go", location, m.parsedData.Model), []byte(content), 0644)

	return err
}
