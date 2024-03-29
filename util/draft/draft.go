package draft

import (
	"github.com/alenn-m/rgen/v2/generator/parser"
)

// Draft represents draft file which is used to generate rest APIs
type Draft struct {
	Models map[string]ModelOptions `yaml:"Models"`
}

// ModelOptions represents single model option
type ModelOptions struct {
	Properties    map[string]string    `yaml:"Properties"`
	Validation    map[string][]string  `yaml:"Validation"`
	Actions       []string             `yaml:"Actions"`
	Relationships parser.Relationships `yaml:"Relationships"`
	OnlyModel     bool                 `yaml:"OnlyModel"`
	Public        bool                 `yaml:"Public"`
}
