package controller

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/files"
	"github.com/alenn-m/rgen/util/templates"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

var dir = "api"

type Input struct {
	Name    string
	Fields  []parser.Field
	Actions []string
}

type Controller struct {
	Input  *Input
	Config *config.Config

	ParsedData parsedData
}

type parsedData struct {
	Root       string
	Package    string
	Controller string
	Model      string
	Fields     string
}

func (c *Controller) Init(input *Input, conf *config.Config) {
	c.Input = input
	c.Config = conf
}

func (c *Controller) Generate() error {
	c.parseData()

	p, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	err = c.createFile(p)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) parseData() {
	c.ParsedData = parsedData{
		Package:    strings.ToLower(inflection.Singular(c.Input.Name)),
		Controller: fmt.Sprintf("%sController", strings.Title(inflection.Plural(c.Input.Name))),
		Model:      strings.Title(inflection.Singular(c.Input.Name)),
		Root:       c.Config.Package,
	}

	for _, item := range c.Input.Fields {
		item.Key = strcase.ToCamel(item.Key)

		c.ParsedData.Fields += fmt.Sprintf("%s: r.%s", item.Key, item.Key) + ",\n"
	}
}

func (c *Controller) createFile(location string) error {
	output, err := ioutil.ReadFile(fmt.Sprintf("%s/src/github.com/alenn-m/rgen/generator/controller/template.tmpl", os.Getenv("GOPATH")))
	if err != nil {
		return err
	}

	content, err := templates.ParseTemplate(string(output), c.ParsedData, map[string]interface{}{
		"ActionUsed": func(input string) bool {
			for _, item := range c.Input.Actions {
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

	servicePath := fmt.Sprintf("%s/%s", location, strings.ToLower(c.Input.Name))
	err = files.MakeDirIfNotExist(servicePath)
	if err != nil {
		return err
	}

	err = c.saveFile([]byte(content), servicePath)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) getServicePath(path string) string {
	return fmt.Sprintf("%s/%s", path, strings.ToLower(c.Input.Name))
}

func (c *Controller) makeDirIfNotExist(location string) error {
	if _, err := os.Stat(location); os.IsNotExist(err) {
		err = os.MkdirAll(location, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) saveFile(content []byte, location string) error {
	err := ioutil.WriteFile(fmt.Sprintf("%s/controller.go", location), content, 0644)

	return err
}
