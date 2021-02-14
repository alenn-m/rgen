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

type Input struct {
	Name    string
	Fields  []parser.Field
	Actions []string
}

type Transport struct {
	Input  *Input
	Config *config.Config

	ParsedData parsedData
}

type parsedData struct {
	Root    string
	Package string
	Prefix  string
	Model   string
	Fields  string
}

func (t *Transport) Init(input *Input, conf *config.Config) {
	t.Input = input
	t.Config = conf
}

func (t *Transport) Generate() error {
	t.parseData()

	content, err := templates.ParseTemplate(TEMPLATE, t.ParsedData, map[string]interface{}{
		"ActionUsed": func(input string) bool {
			for _, item := range t.Input.Actions {
				if item == input {
					return true
				}
			}

			return false
		},
	})
	if err != nil {
		return err
	}

	p, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	err = t.createFile(p, content)
	if err != nil {
		return err
	}

	return nil
}

func (t *Transport) parseData() *Transport {
	t.ParsedData = parsedData{
		Prefix:  strings.ToLower(inflection.Singular(t.Input.Name)),
		Package: strings.ToLower(inflection.Singular(t.Input.Name)),
		Root:    t.Config.Package,
		Model:   strings.Title(inflection.Singular(t.Input.Name)),
	}

	for _, item := range t.Input.Fields {
		t.ParsedData.Fields += fmt.Sprintf("%s %s `json:\"%s\"`\n", strcase.ToCamel(item.Key), item.Value, strcase.ToSnake(item.Key))
	}

	return t
}

func (t *Transport) createFile(location, content string) error {
	servicePath := fmt.Sprintf("%s/%s", location, strings.ToLower(t.Input.Name))

	err := files.MakeDirIfNotExist(servicePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/transport.go", servicePath), []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
