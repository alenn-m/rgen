package generator

import (
	"github.com/alenn-m/rgen/v2/generator/parser"
	"github.com/alenn-m/rgen/v2/util/config"
)

// GeneratorStep represents single generator step
type GeneratorStep interface {
	Generate(input *parser.Parser, conf *config.Config) error
	Save() error
	GetContent() string
}
