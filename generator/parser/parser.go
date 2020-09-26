package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alenn-m/rgen/cmd"
	"github.com/alenn-m/rgen/util/log"
)

type Field struct {
	Key   string
	Value string
}

type Parser struct {
	Name          string
	Fields        []Field
	Actions       []string
	Relationships map[string]string
}

func (p *Parser) Parse(name, fields, actions string) error {
	f := strings.Split(strings.TrimSpace(fields), ",")
	for _, item := range f {
		t := strings.Split(item, ":")
		if len(t) > 1 {
			key, value := t[0], t[1]
			field := Field{
				Key:   key,
				Value: value,
			}

			p.Fields = append(p.Fields, field)
		} else if t[0] == "" {
			return errors.New("fields are required")
		} else {
			return errors.New(fmt.Sprintf("%s has incorrect format", item))
		}
	}

	a := strings.Split(strings.TrimSpace(actions), ",")
	for _, item := range a {
		currentAction := strings.ToUpper(item)
		found := false
		for _, action := range cmd.ACTIONS {
			if action == currentAction {
				found = true
				break
			}
		}

		if !found {
			log.Warning(fmt.Sprintf("Action '%s' is not found, use one of the following [%s]",
				currentAction, strings.Join(cmd.ACTIONS, ", ")))
			continue
		}

		p.Actions = append(p.Actions, currentAction)
	}

	p.Name = strings.TrimSpace(name)

	return nil
}
