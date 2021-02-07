package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
)

func TestModel_Generate(t *testing.T) {
	dir = "test"
	modelName := "Test"

	_, err := os.Stat(fmt.Sprintf("%s/%s", dir, modelName))
	if os.IsExist(err) {
		err = removeTestFile(dir)
		if err != nil {
			t.Error(err.Error())
		}
	}

	m := Model{
		Input: &Input{
			Name: "Test",
			Fields: []parser.Field{
				{
					Key:   "Name",
					Value: "string",
				},
			},
			Relationships: map[string]string{
				"User": "belongsTo",
				"Post": "manyToMany",
			},
		},
		Config:          &config.Config{Package: "github.com/test/testApp"},
		ParsedModelData: parsedModelData{},
	}
	err = m.Generate()
	if err != nil {
		t.Error(err.Error())
	}

	goldenFile, err := ioutil.ReadFile(fmt.Sprintf("%s/Test.golden.go", dir))
	if err != nil {
		t.Error(err.Error())
	}

	createdFile, err := ioutil.ReadFile(fmt.Sprintf("%s/Test.go", dir))
	if err != nil {
		t.Error(err.Error())
	}

	err = m.SetupAutoMigration()
	if err != nil {
		t.Error(err.Error())
	}

	if strings.TrimSpace(string(goldenFile)) != strings.TrimSpace(string(createdFile)) {
		t.Error("Generated file is not correct")
	}

	// cleanup
	err = removeTestFile(dir)
	if err != nil {
		t.Error(err.Error())
	}
}

func removeTestFile(dir string) error {
	err := os.Remove(fmt.Sprintf("%s/Test.go", dir))

	return err
}
