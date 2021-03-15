package draft

import (
	"github.com/alenn-m/rgen/generator/parser"
)

type Draft struct {
	Models map[string]ModelOptions `yaml:"Models"`
}

type ModelOptions struct {
	Properties map[string]string `yaml:"Properties"`
	// Validation    map[string][]string  `yaml:"Validation"`
	Actions       []string             `yaml:"Actions"`
	Relationships parser.Relationships `yaml:"Relationships"`
	OnlyModel     bool                 `yaml:"OnlyModel"`
	Public        bool                 `yaml:"Public"`
}
