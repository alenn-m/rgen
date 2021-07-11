package transport

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/files"
	"github.com/alenn-m/rgen/util/templates"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

//go:embed "template.tmpl"
var TEMPLATE string

var dir = "api"

// parsedData data
type Input struct {
	Name    string
	Fields  []parser.Field
	Actions []string
}

// Transport generator
type Transport struct {
	parsedData parsedData
}

type parsedData struct {
	Root    string
	Package string
	Prefix  string
	Model   string
	Fields  string
	Actions []string
	Content string
}

// Generate generates the 'transport.go' file
func (t *Transport) Generate(input *parser.Parser, conf *config.Config) error {
	t.parseData(input, conf)

	content, err := templates.ParseTemplate(TEMPLATE, t.parsedData, map[string]interface{}{
		"ActionUsed": templates.ActionUsed(t.parsedData.Actions),
	})
	if err != nil {
		return err
	}

	t.parsedData.Content = content

	return nil
}

// GetContent content getter
func (t *Transport) GetContent() string {
	return t.parsedData.Content
}

// Save saves the generated content to file
func (t *Transport) Save() error {
	return t.createFile(t.GetContent())
}

func (t *Transport) parseData(input *parser.Parser, config *config.Config) {
	t.parsedData = parsedData{
		Prefix:  strings.ToLower(inflection.Plural(input.Name)),
		Package: strings.ToLower(inflection.Singular(input.Name)),
		Root:    config.Package,
		Model:   strings.Title(inflection.Singular(input.Name)),
		Actions: input.Actions,
	}

	for _, item := range input.Fields {
		t.parsedData.Fields += fmt.Sprintf("%s %s `json:\"%s\"`\n", strcase.ToCamel(item.Key), item.Value, strcase.ToSnake(item.Key))
	}
}

func (t *Transport) createFile(content string) error {
	location, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	servicePath := fmt.Sprintf("%s/%s", location, t.parsedData.Package)

	err = files.MakeDirIfNotExist(servicePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/transport.go", servicePath), []byte(content), 0644)

	return err
}
