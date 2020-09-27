package service_init

import (
	"fmt"
	"os"
	"testing"

	"github.com/alenn-m/rgen/util/config"
)

func TestServiceInit_Generate(t *testing.T) {
	dir := "test"
	repoName := "KeywordGroup"

	_, err := os.Stat(fmt.Sprintf("%s/%s", dir, repoName))
	if os.IsExist(err) {
		err = removeTestFile(dir)
		if err != nil {
			t.Error(err.Error())
		}
	}

	c := ServiceInit{
		Input: &Input{
			Name: repoName,
		},
		Config: &config.Config{Package: "github.com/test/testApp"},
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
