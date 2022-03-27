package service_init

import (
	"testing"

	"github.com/alenn-m/rgen/v2/generator/parser"
	"github.com/alenn-m/rgen/v2/util/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

var packageName = "github.com/example"

func TestServiceInit_Generate__PrivateRoute(t *testing.T) {
	a := assert.New(t)

	userParser := new(parser.Parser)
	err := userParser.Parse("User", "first_name:string", "")
	a.Nil(err)

	si := new(ServiceInit)
	si.setMainFileLocation("testdata/main.golden")
	err = si.Generate(userParser, &config.Config{Package: packageName})
	a.Nil(err)

	g := goldie.New(t)
	g.Assert(t, "TestServiceInit_Generate__PrivateRoute", []byte(si.GetContent()))
}

func TestServiceInit_Generate__PublicRoute(t *testing.T) {
	a := assert.New(t)

	postParser := new(parser.Parser)
	err := postParser.Parse("Post", "title:string", "")
	a.Nil(err)
	postParser.Public = true

	si := new(ServiceInit)
	si.setMainFileLocation("testdata/main.golden")
	err = si.Generate(postParser, &config.Config{Package: packageName})
	a.Nil(err)

	g := goldie.New(t)
	g.Assert(t, "TestServiceInit_Generate__PublicRoute", []byte(si.GetContent()))
}
