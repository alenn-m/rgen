package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/alenn-m/rgen/v2/generator/parser"
	"github.com/alenn-m/rgen/v2/util/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryImplementation_Generate__Success(t *testing.T) {
	a := assert.New(t)

	p := new(parser.Parser)
	err := p.Parse(modelName, "first_name:string, last_name:string, email:string, age:int", "")
	a.Nil(err)

	p.Relationships = map[string]string{
		"Post":    "hasMany",
		"Profile": "belongsTo",
		"Tag":     "manyToMany",
	}

	repo := &RepositoryImplementation{}
	err = repo.Generate(p, &config.Config{Package: modelName})
	a.Nil(err)

	g := goldie.New(t)
	g.Assert(t, "TestRepositoryImplementation_Generate__Success", []byte(repo.GetContent()))

	err = repo.Save()
	a.Nil(err)

	fp := fmt.Sprintf("%s/%s.go", repo.getServicePath(), repo.parsedData.Package)
	a.FileExists(fp)

	_ = os.RemoveAll(dir)
	_ = os.Remove("api")
}
