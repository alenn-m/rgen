package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alenn-m/rgen/util/misc"
)

// ErrInvalidAction - invalid action
var ErrInvalidAction = fmt.Errorf("Action is not found, use one of the following [%s]",
	strings.Join(misc.ACTIONS, ", "))

// Relationships - list of relationships
type Relationships map[string]string

// Field - single input field
type Field struct {
	Key   string
	Value string
}

// Parser object
type Parser struct {
	Name          string
	Fields        []Field
	Actions       []string
	Relationships Relationships
	OnlyModel     bool
	Public        bool
}

// Parse - parses input data
func (p *Parser) Parse(name, fields, actions string) error {
	p.Name = strings.TrimSpace(name)
	p.Actions = misc.ACTIONS

	f := strings.Split(strings.TrimSpace(fields), ",")
	for _, item := range f {
		t := strings.Split(item, ":")

		if len(t) < 2 {
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
		p.Actions = []string{}

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
				return ErrInvalidAction
			}

			p.Actions = append(p.Actions, currentAction)
		}
	}

	return nil
}
