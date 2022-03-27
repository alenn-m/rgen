package repository

import (
	"fmt"
	"strings"

	"github.com/alenn-m/rgen/v2/generator/parser"
	"github.com/alenn-m/rgen/v2/util/config"
	"github.com/jinzhu/inflection"
)

type parsedData struct {
	Public       bool
	Actions      []string
	Root         string
	Package      string
	Model        string
	Fields       string
	NamedFields  string
	UpdateFields string
	Controller   string
	RepoContent  string
	ImplContent  string
}

func parseData(input *parser.Parser, conf *config.Config) parsedData {
	parsedData := parsedData{
		Public:     input.Public,
		Package:    strings.ToLower(inflection.Singular(input.Name)),
		Model:      strings.Title(inflection.Singular(input.Name)),
		Controller: strings.Title(inflection.Plural(input.Name)) + "Controller",
		Root:       conf.Package,
		Actions:    input.Actions,
	}

	f := []string{}
	nf := []string{}
	uf := []string{}

	for _, item := range input.Fields {
		f = append(f, item.Key)
		nf = append(nf, fmt.Sprintf(":%s", item.Key))
		uf = append(uf, fmt.Sprintf("%s = :%s", item.Key, item.Key))
	}

	parsedData.Fields = strings.Join(f, ", ")
	parsedData.NamedFields = strings.Join(nf, ", ")
	parsedData.UpdateFields = strings.Join(uf, ", ")

	return parsedData
}
