package controller

import (
	"path/filepath"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
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
	Package    string
	Controller string
	Model      string
	Fields     string
}

func (c *Controller) Init(input *Input, conf *config.Config) {
	// TODO: parse actions into array of strings, for now all actions are created
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
