package controller

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

func TestController_Generate__Success(t *testing.T) {
	a := assert.New(t)

	p := new(parser.Parser)
	p.Parse(modelName, "first_name:string, last_name:string, email:string, age:int", "")
	p.Relationships = map[string]string{
		"Post":    "hasMany",
		"Profile": "belongsTo",
		"Tag":     "manyToMany",
	}

	c := &Controller{}
	err := c.Generate(p, &config.Config{Package: modelName})
	a.Nil(err)

	g := goldie.New(t)
	g.Assert(t, "TestController_Generate__Success", []byte(c.GetContent()))

	err = c.Save()
	a.Nil(err)

	fp := fmt.Sprintf("%s/%s/controller.go", dir, c.parsedData.Package)
	a.FileExists(fp)

	_ = os.RemoveAll(dir)
}
