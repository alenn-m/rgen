package repository

import (
	_ "embed"
	"fmt"
	"go/format"
	"io/ioutil"

	"github.com/alenn-m/rgen/v2/generator/parser"
	"github.com/alenn-m/rgen/v2/util/config"
	"github.com/alenn-m/rgen/v2/util/files"
	"github.com/alenn-m/rgen/v2/util/templates"
	"github.com/jinzhu/inflection"
)

//go:embed "template_mysql.tmpl"
var TEMPLATE_MYSQL string

// RepositoryImplementation generator
type RepositoryImplementation struct {
	parsedData parsedData
}

// Generate generates repository_imlementation
func (r *RepositoryImplementation) Generate(input *parser.Parser, conf *config.Config) error {
	r.parsedData = parseData(input, conf)

	contentString, err := templates.ParseTemplate(TEMPLATE_MYSQL, r.parsedData, map[string]interface{}{
		"ActionUsed": templates.ActionUsed(r.parsedData.Actions),
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

// Save saves generated repository_implementation to file
func (r *RepositoryImplementation) Save() error {
	repositoriesPath := r.getServicePath()
	err := files.MakeDirIfNotExist(repositoriesPath)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/%s.go", repositoriesPath, r.parsedData.Package), []byte(r.GetContent()), 0644)
}

// GetContent returns generated repository_implementation to file
func (r *RepositoryImplementation) GetContent() string {
	return r.parsedData.ImplContent
}

func (r *RepositoryImplementation) getServicePath() string {
	return fmt.Sprintf("%s/%s/repositories/mysql", dir, r.parsedData.Package)
}
