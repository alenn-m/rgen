package controller

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/files"
	"github.com/alenn-m/rgen/util/templates"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

//go:embed "template.tmpl"
var template string

var dir = "api"

type Controller struct {
	parsedData parsedData
}

func (c *Controller) Generate(input *parser.Parser, conf *config.Config) error {
	c.parseData(input, conf)

	content, err := templates.ParseTemplate(template, c.parsedData, map[string]interface{}{
		"ActionUsed": func(input string) bool {
			for _, item := range c.parsedData.Actions {
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

	c.parsedData.Content = content

	return nil
}

func (c *Controller) Save() error {
	servicePath := fmt.Sprintf("%s/%s", dir, strings.ToLower(c.parsedData.Name))
	err := files.MakeDirIfNotExist(servicePath)
	if err != nil {
		return err
	}

	return c.saveFile([]byte(c.GetContent()), servicePath)
}

func (c *Controller) GetContent() string {
	return c.parsedData.Content
}

type parsedData struct {
	Name       string
	Actions    []string
	Root       string
	Package    string
	Controller string
	Model      string
	Fields     string
	Content    string
}

func (c *Controller) parseData(input *parser.Parser, conf *config.Config) {
	c.parsedData = parsedData{
		Name:       input.Name,
		Actions:    input.Actions,
		Package:    strings.ToLower(inflection.Singular(input.Name)),
		Controller: fmt.Sprintf("%sController", strings.Title(inflection.Plural(input.Name))),
		Model:      strings.Title(inflection.Singular(input.Name)),
		Root:       conf.Package,
	}

	for _, item := range input.Fields {
		item.Key = strcase.ToCamel(item.Key)

		c.parsedData.Fields += fmt.Sprintf("%s: r.%s", item.Key, item.Key) + ",\n"
	}
}

func (c *Controller) saveFile(content []byte, location string) error {
	err := ioutil.WriteFile(fmt.Sprintf("%s/controller.go", location), content, 0644)

	return err
}
