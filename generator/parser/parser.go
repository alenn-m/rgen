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
	SkipController bool
	Actions        []string
}

func (p *Parser) Parse(name, fields, rootDir, actions string) error {
	// Note: currently unused
	p.RootDir = rootDir

	splitActions := strings.Split(actions, ",")
	for _, item := range splitActions {
		p.Actions = append(p.Actions, strings.TrimSpace(item))
	}

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
