package repository

import (
	_ "embed"
	"fmt"
	"go/format"
	"io/ioutil"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/files"
	"github.com/alenn-m/rgen/util/templates"
	"github.com/jinzhu/inflection"
)

//go:embed "template_mysql.tmpl"
var TEMPLATE_MYSQL string

type RepositoryImplementation struct {
	parsedData parsedData
}

func (r *RepositoryImplementation) Generate(input *parser.Parser, conf *config.Config) error {
	r.parsedData = parseData(input, conf)

	contentString, err := templates.ParseTemplate(TEMPLATE_MYSQL, r.parsedData, map[string]interface{}{
		"ActionUsed": func(input string) bool {
			for _, item := range r.parsedData.Actions {
				if item == input {
					return true
				}
			}

			return false
		},
		"Pluralize": func(input string) string {
			return inflection.Plural(input)
		},
	})
	if err != nil {
		return err
	}

	mysqlContent, err := format.Source([]byte(contentString))
	if err != nil {
		return err
	}

	r.parsedData.ImplContent = string(mysqlContent)

	return nil
}

func (r *RepositoryImplementation) Save() error {
	repositoriesPath := r.getServicePath(dir)
	err := files.MakeDirIfNotExist(repositoriesPath)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/%s.go", repositoriesPath, r.parsedData.Package), []byte(r.GetContent()), 0644)
}

func (r *RepositoryImplementation) GetContent() string {
	return r.parsedData.ImplContent
}

func (r *RepositoryImplementation) getServicePath(servicePath string) string {
	return fmt.Sprintf("%s/repositories/mysql", servicePath)
}
