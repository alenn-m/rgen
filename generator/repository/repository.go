package repository

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/misc"
	"github.com/jinzhu/inflection"
)

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

	mysqlActions := []string{MYSQL_TEMPLATE_HEADER, MYSQL_TEMPLATE_INDEX, MYSQL_TEMPLATE_SHOW, MYSQL_TEMPLATE_CREATE, MYSQL_TEMPLATE_UPDATE, MYSQL_TEMPLATE_DELETE}
	dbRepositoryActions := []string{DBR_INDEX, DBR_SHOW, DBR_CREATE, DBR_UPDATE, DBR_DELETE}
	repositoryActions := []string{R_INDEX, R_SHOW, R_CREATE, R_UPDATE, R_DELETE}

	if len(r.Input.Actions) > 0 {
		mysqlActions = []string{MYSQL_TEMPLATE_HEADER}
		dbRepositoryActions = []string{}
		repositoryActions = []string{}

		for _, action := range r.Input.Actions {
			switch action {
			case misc.ACTION_INDEX:
				mysqlActions = append(mysqlActions, MYSQL_TEMPLATE_INDEX)
				dbRepositoryActions = append(dbRepositoryActions, DBR_INDEX)
				repositoryActions = append(repositoryActions, R_INDEX)
			case misc.ACTION_SHOW:
				mysqlActions = append(mysqlActions, MYSQL_TEMPLATE_SHOW)
				dbRepositoryActions = append(dbRepositoryActions, DBR_SHOW)
				repositoryActions = append(repositoryActions, R_SHOW)
			case misc.ACTION_CREATE:
				mysqlActions = append(mysqlActions, MYSQL_TEMPLATE_CREATE)
				dbRepositoryActions = append(dbRepositoryActions, DBR_CREATE)
				repositoryActions = append(repositoryActions, R_CREATE)
			case misc.ACTION_UPDATE:
				mysqlActions = append(mysqlActions, MYSQL_TEMPLATE_UPDATE)
				dbRepositoryActions = append(dbRepositoryActions, DBR_UPDATE)
				repositoryActions = append(repositoryActions, R_UPDATE)
			case misc.ACTION_DELETE:
				mysqlActions = append(mysqlActions, MYSQL_TEMPLATE_DELETE)
				dbRepositoryActions = append(dbRepositoryActions, DBR_DELETE)
				repositoryActions = append(repositoryActions, R_DELETE)
			}
		}
	}

	contentString := TEMPLATE_AUTH
	if r.Input.Public {
		contentString = TEMPLATE_NO_AUTH
	}
	contentString = strings.Replace(contentString, "{{DBRepositoryActions}}", strings.Join(dbRepositoryActions, "\n"), -1)
	contentString = strings.Replace(contentString, "{{RepositoryActions}}", strings.Join(repositoryActions, "\n"), -1)
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

	contentString = strings.Join(mysqlActions, "\n")
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
	r.ParsedData.Model = strings.Title(inflection.Singular(r.Input.Name))
}

func (r *Repository) parsePackage() {
	r.ParsedData.Package = strings.ToLower(inflection.Singular(r.Input.Name))
}

func (r *Repository) parseController() {
	r.ParsedData.Controller = strings.Title(inflection.Plural(r.Input.Name)) + "Controller"
}

func (r *Repository) saveFile(content []byte, location string) error {
	err := ioutil.WriteFile(fmt.Sprintf("%s/repository.go", location), content, 0644)

	return err
}
