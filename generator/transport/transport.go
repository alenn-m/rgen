package transport

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/misc"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

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
	t.
		parsePackage().
		parseModelName().
		parsePrefix().
		parseFields()

	p, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	err = t.createFile(p)
	if err != nil {
		return err
	}

	return nil
}

func (t *Transport) parsePrefix() *Transport {
	t.ParsedData.Prefix = strings.ToLower(inflection.Singular(t.Input.Name))

	return t
}

func (t *Transport) parsePackage() *Transport {
	t.ParsedData.Package = strings.ToLower(inflection.Singular(t.Input.Name))

	return t
}

func (t *Transport) parseModelName() *Transport {
	t.ParsedData.Model = strings.Title(inflection.Singular(t.Input.Name))

	return t
}

func (t *Transport) parseFields() *Transport {
	t.ParsedData.Fields = ""
	for _, item := range t.Input.Fields {
		t.ParsedData.Fields += fmt.Sprintf("%s %s `json:\"%s\"`\n", strcase.ToCamel(item.Key), item.Value, strcase.ToSnake(item.Key))
	}

	return t
}

func (t *Transport) createFile(location string) error {
	actions := []string{TRANSPORT_HEADER, TRANSPORT_INDEX, TRANSPORT_SHOW, TRANSPORT_CREATE, TRANSPORT_UPDATE, TRANSPORT_DELETE}
	tActions := []string{T_INDEX, T_SHOW, T_CREATE, T_UPDATE, T_DELETE}

	if len(t.Input.Actions) > 0 {
		actions = []string{TRANSPORT_HEADER}
		tActions = []string{}

		for _, action := range t.Input.Actions {
			switch action {
			case misc.ACTION_INDEX:
				actions = append(actions, TRANSPORT_INDEX)
				tActions = append(tActions, T_INDEX)
			case misc.ACTION_SHOW:
				actions = append(actions, TRANSPORT_SHOW)
				tActions = append(tActions, T_SHOW)
			case misc.ACTION_CREATE:
				actions = append(actions, TRANSPORT_CREATE)
				tActions = append(tActions, T_CREATE)
			case misc.ACTION_UPDATE:
				actions = append(actions, TRANSPORT_UPDATE)
				tActions = append(tActions, T_UPDATE)
			case misc.ACTION_DELETE:
				actions = append(actions, TRANSPORT_DELETE)
				tActions = append(tActions, T_DELETE)
			}
		}
	}

	contentString := strings.Join(actions, "\n")
	contentString = strings.Replace(contentString, "{{TransportActions}}", strings.Join(tActions, "\n"), -1)
	contentString = strings.Replace(contentString, "{{Package}}", t.ParsedData.Package, -1)
	contentString = strings.Replace(contentString, "{{Prefix}}", t.ParsedData.Prefix, -1)
	contentString = strings.Replace(contentString, "{{Model}}", t.ParsedData.Model, -1)
	contentString = strings.Replace(contentString, "{{Fields}}", t.ParsedData.Fields, -1)
	contentString = strings.Replace(contentString, "{{Root}}", t.Config.Package, -1)

	content, err := format.Source([]byte(contentString))
	if err != nil {
		return err
	}

	servicePath := t.getServicePath(location)
	err = t.makeDirIfNotExist(servicePath)
	if err != nil {
		return err
	}

	err = t.saveFile(content, servicePath)
	if err != nil {
		return err
	}

	return nil
}

func (t *Transport) getServicePath(path string) string {
	return fmt.Sprintf("%s/%s", path, strings.ToLower(t.Input.Name))
}

func (t *Transport) makeDirIfNotExist(location string) error {
	if _, err := os.Stat(location); os.IsNotExist(err) {
		err = os.MkdirAll(location, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Transport) saveFile(content []byte, location string) error {
	err := ioutil.WriteFile(fmt.Sprintf("%s/transport.go", location), content, 0644)

	return err
}
