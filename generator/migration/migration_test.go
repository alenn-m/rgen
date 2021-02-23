package migration

import (
	"testing"

	"github.com/alenn-m/rgen/generator/parser"
)

func TestMigration_Generate(t *testing.T) {
	m := Migration{
		Input: &Input{
			Name: "Post",
			Fields: []parser.Field{
				{Key: "Title", Value: "string"},
				{Key: "Body", Value: "string"},
				{Key: "UserID", Value: "int64"},
			},
		},
		Config: nil,
	}
	err := m.Generate()
	if err != nil {
		t.Error(err.Error())
	}
}
