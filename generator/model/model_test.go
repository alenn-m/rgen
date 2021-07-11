package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

var modelName = "User"

func TestModel_Generate__Success(t *testing.T) {
	a := assert.New(t)

	p := new(parser.Parser)
	err := p.Parse(modelName, "first_name:string, last_name:string, email:string, age:int", "")
	a.Nil(err)

	p.Relationships = map[string]string{
		"Post":    "hasMany",
		"Profile": "belongsTo",
		"Tag":     "manyToMany",
	}

	model := &Model{}
	err = model.Generate(p, &config.Config{Package: modelName})
	a.Nil(err)

	g := goldie.New(t)
	g.Assert(t, "TestModel_Generate__Success", []byte(model.GetContent()))

	err = model.Save()
	a.Nil(err)

	fp := fmt.Sprintf("%s/%s.go", dir, modelName)
	a.FileExists(fp)

	_ = os.RemoveAll(dir)
}

func TestModel_Generate__WrongRelationship(t *testing.T) {
	a := assert.New(t)

	p := new(parser.Parser)
	err := p.Parse(modelName, "first_name:string, last_name:string, email:string, age:int", "")
	a.Nil(err)

	p.Relationships = map[string]string{
		"Post":    "hasMany",
		"Profile": "wrongRelationship",
	}

	model := &Model{}
	err = model.Generate(p, &config.Config{Package: modelName})
	a.NotNil(err)
	a.Equal(err, ErrInvalidRelationship)
}
