package migration

import (
	"testing"
)

func TestPivotMigration_Generate(t *testing.T) {
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
	pm.Init(tables)

	err := pm.Generate()
	if err != nil {
		t.Errorf(err.Error())
	}
}
