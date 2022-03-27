package transport

import (
	"fmt"
	"os"
	"testing"

	"github.com/alenn-m/rgen/v2/generator/parser"
	"github.com/alenn-m/rgen/v2/util/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

var modelName = "User"

func TestTransport_Generate(t *testing.T) {
	a := assert.New(t)

	p := new(parser.Parser)
	p.Validation = parser.Validation{
		"name": []string{"Required", "Min:5"},
	}
	err := p.Parse(modelName, "name:string", "")
	a.Nil(err)

	c := Transport{}
	err = c.Generate(p, &config.Config{Package: "github.com/test/testApp"})
	a.Nil(err)

	g := goldie.New(t)
	g.Assert(t, "TestTransport_Generate", []byte(c.GetContent()))

	err = c.Save()
	a.Nil(err)

	fp := fmt.Sprintf("%s/%s/transport.go", dir, c.parsedData.Package)
	a.FileExists(fp)

	_ = os.RemoveAll(dir)
}
