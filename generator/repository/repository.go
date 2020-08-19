package repository

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alenn-m/rgen/util/config"
	"github.com/jinzhu/inflection"
)

var dir = "api"

type Input struct {
	Name string
}

type Repository struct {
	Input      *Input
	Config     *config.Config
	ParsedData parsedData
}

type parsedData struct {
	Package    string
	Model      string
	Controller string
}

func (r *Repository) Init(input *Input, conf *config.Config) {
	r.Input = input
	r.Config = conf
}

func (r *Repository) Generate() error {
	r.parseModel()
	r.parseController()
	r.parsePackage()

	contentString := TEMPLATE
	contentString = strings.Replace(contentString, "{{Model}}", r.ParsedData.Model, -1)
	contentString = strings.Replace(contentString, "{{Package}}", r.ParsedData.Package, -1)
	contentString = strings.Replace(contentString, "{{Controller}}", r.ParsedData.Controller, -1)
	contentString = strings.Replace(contentString, "{{Root}}", r.Config.Package, -1)

	content, err := format.Source([]byte(contentString))
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

	err = r.saveFile(content, servicePath)
	if err != nil {
		return err
	}

	contentString = MYSQL_TEMPLATE
	contentString = strings.Replace(contentString, "{{Model}}", r.ParsedData.Model, -1)
	contentString = strings.Replace(contentString, "{{Root}}", r.Config.Package, -1)

	content, err = format.Source([]byte(contentString))
	if err != nil {
		return err
	}

	repositoriesPath := r.getRepositoryPath(servicePath)
	err = r.makeDirIfNotExist(repositoriesPath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.go", repositoriesPath, r.ParsedData.Package), content, 0644)
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

func (r *Repository) parseModel() {
	r.ParsedData.Model = strings.Title(strings.ToLower(inflection.Singular(r.Input.Name)))
}

func (r *Repository) parsePackage() {
	r.ParsedData.Package = strings.ToLower(inflection.Singular(r.Input.Name))
}

func (r *Repository) parseController() {
	r.ParsedData.Controller = strings.Title(strings.ToLower(inflection.Plural(r.Input.Name))) + "Controller"
}

func (r *Repository) saveFile(content []byte, location string) error {
	err := ioutil.WriteFile(fmt.Sprintf("%s/repository.go", location), content, 0644)

	return err
}
