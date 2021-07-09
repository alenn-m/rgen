package migration

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

var modelName = "User"

func TestMigration_Generate__Success(t *testing.T) {
	a := assert.New(t)

	p := new(parser.Parser)
	p.Parse(modelName, "first_name:string, last_name:string, email:string, age:int", "")
	p.Relationships = map[string]string{
		"Post":    "hasMany",
		"Profile": "belongsTo",
		"Tag":     "manyToMany",
	}

	m := &Migration{}
	err := m.Generate(p, &config.Config{Package: modelName})
	a.Nil(err)

	g := goldie.New(t)
	g.Assert(t, "TestMigration_Generate__Success", []byte(m.GetContent()))

	err = m.Save()
	a.Nil(err)

	files, err := ioutil.ReadDir(dir)
	a.Nil(err)

	a.Equal(1, len(files))

	file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, files[0].Name()))
	a.Nil(err)

	a.Equal(m.GetContent(), string(file))

	err = os.RemoveAll(dir)
	a.Nil(err)

	err = os.Remove("database")
	a.Nil(err)
}
