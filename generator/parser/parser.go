package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alenn-m/rgen/v2/util/misc"
)

// ErrInvalidAction - invalid action
var ErrInvalidAction = fmt.Errorf("Action is not found, use one of the following [%s]",
	strings.Join(misc.ACTIONS, ", "))

// Relationships - list of relationships
type Relationships map[string]string

// Validation - list of validations
type Validation map[string][]string

// Field - single input field
// Field represents a single input field. It has a Key and a Value.
type Field struct {
	Key   string
	Value string
}

// Parser object
// Parser represents a parser object. It has a Name, Fields, Actions, Relationships, Validation, OnlyModel, and Public.
type Parser struct {
	Name          string
	Fields        []Field
	Actions       []string
	Relationships Relationships
	Validation    Validation
	OnlyModel     bool
	Public        bool
}

// Parse - parses input data
// Parse takes in a name, fields, and actions as inputs, and parses them into a Parser object.
func (p *Parser) Parse(name, fields, actions string) error {
	p.Name = strings.TrimSpace(name)
	p.Actions = misc.ACTIONS
	p.Validation = make(Validation)

	f := strings.Split(strings.TrimSpace(fields), ",")
	for _, item := range f {
		t := strings.SplitN(item, ":", 2)

		if len(t) < 2 {
			return fmt.Errorf("%s has incorrect format", item)
		}
		if len(t) > 1 {
			key, value := strings.TrimSpace(t[0]), strings.TrimSpace(t[1])

			r := strings.Split(value, "#")

			field := Field{
				Key:   key,
				Value: r[0],
			}

			if len(r) > 1 {
				validations := strings.Split(r[1], "|")
				for _, v := range validations {
					if _, found := p.Validation[key]; found {
						p.Validation[key] = append(p.Validation[key], strings.TrimSpace(v))
					} else {
						p.Validation[key] = []string{strings.TrimSpace(v)}
					}
				}
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
