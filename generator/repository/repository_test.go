package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/misc"
)

func TestRepository_Generate(t *testing.T) {
	dir = "test"
	repoName := "Test"

	_, err := os.Stat(fmt.Sprintf("%s/%s", dir, repoName))
	if os.IsExist(err) {
		err = removeTestFile(dir)
		if err != nil {
			t.Error(err.Error())
		}
	}

	c := Repository{
		Input: &Input{
			Name:    repoName,
			Actions: []string{misc.ACTION_CREATE, misc.ACTION_UPDATE},
		},
		Config:     &config.Config{Package: "github.com/test/testApp"},
		ParsedData: parsedData{},
	}

	err = c.Generate()
	if err != nil {
		t.Error(err.Error())
	}
}

func removeTestFile(dir string) error {
	err := os.Remove(fmt.Sprintf("%s/Test.go", dir))

	return err
}
