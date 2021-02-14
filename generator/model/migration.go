package model

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alenn-m/rgen/util/files"
	"github.com/alenn-m/rgen/util/templates"
)

//go:embed "migrations_template.tmpl"
var MIGRATIONS_TEMPLATE string

func (m *Model) SetupAutoMigration() error {
	location, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	err = m.parseModels(location)
	if err != nil {
		return err
	}

	// in case someone removed this file, let's create it again
	if !files.FileExists("migrations.go") {
		err = ioutil.WriteFile("migrations.go", []byte(""), 0644)
		if err != nil {
			return err
		}
	}

	migrationsLocation, err := filepath.Abs("./")
	if err != nil {
		return err
	}

	content, err := templates.ParseTemplate(MIGRATIONS_TEMPLATE, m.ParsedMigrationData, nil)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", migrationsLocation, "migrations.go"), []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) parseModels(location string) error {
	models := []string{}

	err := filepath.Walk(location, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			re := regexp.MustCompile("(?m)^.*?(type\\b.*?(struct)[\\s]*)")
			match := re.FindString(string(data))
			if len(match) > 0 {
				match = strings.Replace(match, "type", "", -1)
				match = strings.Replace(match, "struct", "", -1)
				modelName := strings.TrimSpace(match)
				if modelName == "Base" {
					return nil
				}

				models = append(models, fmt.Sprintf("&models.%s{}", modelName))
			}
			return err
		}

		return nil
	})

	m.ParsedMigrationData = parsedMigrationData{
		Models: strings.Join(models, ", "),
		Root:   m.Config.Package,
	}

	return err
}
