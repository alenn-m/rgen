package transport

import (
	"fmt"
	"os"
	"testing"

	"github.com/alenn-m/rgen/util/config"
	"github.com/stretchr/testify/assert"
)

func TestTransport_Generate(t *testing.T) {
	a := assert.New(t)

	dir = "test"
	repoName := "Test"

	_, err := os.Stat(fmt.Sprintf("%s/%s", dir, repoName))
	if os.IsExist(err) {
		err = removeTestFile(dir)
		a.Nil(err)
	}

	c := Transport{
		Input: &Input{
			Name:    repoName,
			Actions: []string{},
		},
		Config:     &config.Config{Package: "github.com/test/testApp"},
		ParsedData: parsedData{},
	}

	err = c.Generate()
	a.Nil(err)

	a.Equal("github.com/test/testApp", c.ParsedData.Root)
	a.Equal("test", c.ParsedData.Package)
	a.Equal("Test", c.ParsedData.Model)
	a.Equal("test", c.ParsedData.Prefix)
}

func removeTestFile(dir string) error {
	err := os.Remove(fmt.Sprintf("%s/Test.go", dir))

	return err
}
