package migration

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/alenn-m/rgen/v2/util/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestPivotMigration_Generate(t *testing.T) {
	a := assert.New(t)

	pm := new(PivotMigration)

	tables := []PivotMigrationEntry{
		{
			TableOne: "User",
			TableTwo: "Post",
		},
		{
			TableOne: "User",
			TableTwo: "Comment",
		},
		{
			TableOne: "User",
			TableTwo: "Post",
		},
	}
	pm.Init(tables, &config.Config{
		Migration: config.Migration{
			Sequential: true,
		},
	})

	err := pm.Generate()
	a.Nil(err)

	files, err := ioutil.ReadDir(dir)
	a.Nil(err)
	a.Equal(2, len(files))

	g := goldie.New(t)

	for _, file := range files {
		f, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, file.Name()))
		a.Nil(err)

		g.Assert(t, fmt.Sprintf("TestPivotMigration_Generate_%s", file.Name()), f)
	}

	a.Nil(os.RemoveAll(dir))
	a.Nil(os.Remove("database"))
}
