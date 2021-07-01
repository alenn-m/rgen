package generator

import (
	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/config"
)

// GeneratorStep represents single generator step
type GeneratorStep interface {
	Generate(input *parser.Parser, conf *config.Config) error
	Save() error
	GetContent() string
}
