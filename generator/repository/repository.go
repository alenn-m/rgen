package repository

import (
	_ "embed"
	"fmt"
	"io/ioutil"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/files"
	"github.com/alenn-m/rgen/util/templates"
)

//go:embed "template_auth.tmpl"
var TEMPLATE_AUTH string

//go:embed "template_no_auth.tmpl"
var TEMPLATE_NO_AUTH string

var dir = "api"

// Repository generator
type Repository struct {
	parsedData parsedData
}

// Generate generates repository
func (r *Repository) Generate(input *parser.Parser, conf *config.Config) error {
	r.parsedData = parseData(input, conf)

	contentString := TEMPLATE_AUTH
	if r.parsedData.Public {
		contentString = TEMPLATE_NO_AUTH
	}

	content, err := templates.ParseTemplate(contentString, r.parsedData, map[string]interface{}{
		"ActionUsed": func(input string) bool {
			for _, item := range r.parsedData.Actions {
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

	r.parsedData.RepoContent = content

	return nil
}

// Save saves repository to file
func (r *Repository) Save() error {
	servicePath := r.getServicePath(dir)
	err := files.MakeDirIfNotExist(servicePath)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/repository.go", servicePath), []byte(r.GetContent()), 0644)
}

// GetContent returns generated repository content
func (r *Repository) GetContent() string {
	return r.parsedData.RepoContent
}

func (r *Repository) getServicePath(path string) string {
	return fmt.Sprintf("%s/%s", path, r.parsedData.Package)
}
