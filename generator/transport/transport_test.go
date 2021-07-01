package transport

import (
	"testing"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestTransport_Generate(t *testing.T) {
	a := assert.New(t)

	dir = "test"
	repoName := "Test"

	p := new(parser.Parser)
	p.Parse(repoName, "", "")

	c := Transport{}
	err := c.Generate(p, &config.Config{Package: "github.com/test/testApp"})
	a.Nil(err)

	g := goldie.New(t)
	g.Assert(t, "TestTransport_Generate", []byte(c.GetContent()))
}
