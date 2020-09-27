package controller

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

func (c *Controller) parseData() {
	c.ParsedData = parsedData{}
	c.
		parsePackage().
		parseControllerName().
		parseModelName().
		parseFields()
}

func (c *Controller) parsePackage() *Controller {
	c.ParsedData.Package = strings.ToLower(inflection.Singular(c.Input.Name))

	return c
}

func (c *Controller) parseControllerName() *Controller {
	c.ParsedData.Controller = fmt.Sprintf("%sController", strings.Title(inflection.Plural(c.Input.Name)))

	return c
}

func (c *Controller) parseModelName() *Controller {
	c.ParsedData.Model = strings.Title(inflection.Singular(c.Input.Name))

	return c
}

func (c *Controller) parseFields() *Controller {
	c.ParsedData.Fields = ""
	for _, item := range c.Input.Fields {
		item.Key = strcase.ToCamel(item.Key)

		c.ParsedData.Fields += fmt.Sprintf("%s: r.%s", item.Key, item.Key) + ",\n"
	}

	return c
}
