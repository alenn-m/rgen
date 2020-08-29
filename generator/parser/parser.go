package parser

import (
	"errors"
	"fmt"
	"strings"
)

type Field struct {
	Key   string
	Value string
}

type Parser struct {
	Name           string
	Fields         []Field
	Relationships  map[string]string
	RootDir        string
	MakeController bool
}

func (p *Parser) Parse(name, fields, rootDir string) error {
	// Note: currently unused
	p.RootDir = rootDir

	fs := []Field{}

	f := strings.Split(strings.TrimSpace(fields), ",")
	for _, item := range f {
		t := strings.Split(item, ":")
		if len(t) > 1 {
			key, value := t[0], t[1]
			field := Field{
				Key:   key,
				Value: value,
			}

			fs = append(fs, field)
		} else if t[0] == "" {
			return errors.New("fields are required")
		} else {
			return errors.New(fmt.Sprintf("%s has incorrect format", item))
		}
	}

	p.Fields = fs
	p.Name = strings.TrimSpace(name)

	return nil
}
