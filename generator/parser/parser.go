package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alenn-m/rgen/util/log"
	"github.com/alenn-m/rgen/util/misc"
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
	OnlyModel     bool
	Public        bool
}

func (p *Parser) Parse(name, fields, actions string) error {
	f := strings.Split(strings.TrimSpace(fields), ",")
	for _, item := range f {
		t := strings.Split(item, ":")
		if len(t) < 1{
			return fmt.Errorf("%s has incorrect format", item)
		}
		if len(t) > 1 {
			key, value := t[0], t[1]
			field := Field{
				Key:   key,
				Value: value,
			}

			p.Fields = append(p.Fields, field)
		} else if t[0] == "" {
			return errors.New("fields are required")
		}
	}

	if actions != "" {
		a := strings.Split(strings.TrimSpace(actions), ",")
		for _, item := range a {
			currentAction := strings.ToUpper(strings.TrimSpace(item))
			found := false
			for _, action := range misc.ACTIONS {
				if action == currentAction {
					found = true
					break
				}
			}

			if !found {
				log.Warning(fmt.Sprintf("Action '%s' is not found, use one of the following [%s]",
					currentAction, strings.Join(misc.ACTIONS, ", ")))
				continue
			}

			p.Actions = append(p.Actions, currentAction)
		}
	} else {
		p.Actions = misc.ACTIONS
	}

	p.Name = strings.TrimSpace(name)

	return nil
}
