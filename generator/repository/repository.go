package repository

import (
	_ "embed"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/templates"
	"github.com/jinzhu/inflection"
)

//go:embed "template_auth.tmpl"
var TEMPLATE_AUTH string

//go:embed "template_no_auth.tmpl"
var TEMPLATE_NO_AUTH string

//go:embed "template_mysql.tmpl"
var TEMPLATE_MYSQL string

var dir = "api"

type Input struct {
	Name    string
	Actions []string
	Public  bool
}

type Repository struct {
	Input      *Input
	Config     *config.Config
	ParsedData parsedData
}

type parsedData struct {
	Root       string
	Package    string
	Model      string
	Controller string
}

func (r *Repository) Init(input *Input, conf *config.Config) {
	r.Input = input
	r.Config = conf
}

func (r *Repository) Generate() error {
	r.parseData()

	contentString := TEMPLATE_AUTH
	if r.Input.Public {
		contentString = TEMPLATE_NO_AUTH
	}

	content, err := templates.ParseTemplate(contentString, r.ParsedData, map[string]interface{}{
		"ActionUsed": func(input string) bool {
			for _, item := range r.Input.Actions {
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

	location, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	servicePath := r.getServicePath(location)
	err = r.makeDirIfNotExist(servicePath)
	if err != nil {
		return err
	}

	err = r.saveFile([]byte(content), servicePath)
	if err != nil {
		return err
	}

	contentString, err = templates.ParseTemplate(TEMPLATE_MYSQL, r.ParsedData, map[string]interface{}{
		"ActionUsed": func(input string) bool {
			for _, item := range r.Input.Actions {
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

	mysqlContent, err := format.Source([]byte(contentString))
	if err != nil {
		return err
	}

	repositoriesPath := r.getRepositoryPath(servicePath)
	err = r.makeDirIfNotExist(repositoriesPath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.go", repositoriesPath, r.ParsedData.Package), mysqlContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) getServicePath(path string) string {
	return fmt.Sprintf("%s/%s", path, r.ParsedData.Package)
}

func (r *Repository) getRepositoryPath(servicePath string) string {
	return fmt.Sprintf("%s/repositories/mysql", servicePath)
}

func (r *Repository) makeDirIfNotExist(location string) error {
	if _, err := os.Stat(location); os.IsNotExist(err) {
		err = os.MkdirAll(location, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) parseData() {
	r.ParsedData = parsedData{
		Package:    strings.ToLower(inflection.Singular(r.Input.Name)),
		Model:      strings.Title(inflection.Singular(r.Input.Name)),
		Controller: strings.Title(inflection.Plural(r.Input.Name)) + "Controller",
		Root:       r.Config.Package,
	}
}

func (r *Repository) saveFile(content []byte, location string) error {
	err := ioutil.WriteFile(fmt.Sprintf("%s/repository.go", location), content, 0644)

	return err
}
