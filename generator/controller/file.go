package controller

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alenn-m/rgen/util/misc"
)

func (c *Controller) createFile(location string) error {
	actions := []string{CONTROLLER_HEADER, CONTROLLER_INDEX, CONTROLLER_STORE, CONTROLLER_UPDATE, CONTROLLER_DELETE}

	if len(c.Input.Actions) > 0 {
		actions = []string{CONTROLLER_HEADER}

		for _, action := range c.Input.Actions {
			switch action {
			case misc.ACTION_READ:
				actions = append(actions, CONTROLLER_INDEX)
			case misc.ACTION_CREATE:
				actions = append(actions, CONTROLLER_STORE)
			case misc.ACTION_UPDATE:
				actions = append(actions, CONTROLLER_UPDATE)
			case misc.ACTION_DELETE:
				actions = append(actions, CONTROLLER_DELETE)
			}
		}
	}

	contentString := strings.Join(actions, "\n")
	contentString = strings.Replace(contentString, "{{Package}}", c.ParsedData.Package, -1)
	contentString = strings.Replace(contentString, "{{Controller}}", c.ParsedData.Controller, -1)
	contentString = strings.Replace(contentString, "{{Model}}", c.ParsedData.Model, -1)
	contentString = strings.Replace(contentString, "{{Fields}}", c.ParsedData.Fields, -1)
	contentString = strings.Replace(contentString, "{{Root}}", c.Config.Package, -1)

	content, err := format.Source([]byte(contentString))
	if err != nil {
		return err
	}

	servicePath := c.getServicePath(location)
	err = c.makeDirIfNotExist(servicePath)
	if err != nil {
		return err
	}

	err = c.saveFile(content, servicePath)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) getServicePath(path string) string {
	return fmt.Sprintf("%s/%s", path, strings.ToLower(c.Input.Name))
}

func (c *Controller) makeDirIfNotExist(location string) error {
	if _, err := os.Stat(location); os.IsNotExist(err) {
		err = os.MkdirAll(location, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) saveFile(content []byte, location string) error {
	err := ioutil.WriteFile(fmt.Sprintf("%s/controller.go", location), content, 0644)

	return err
}
