package model

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alenn-m/rgen/util/files"
)

func (m *Model) SetupAutoMigration() error {
	location, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	models := []string{}

	err = filepath.Walk(location, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			re := regexp.MustCompile("type([^|]*)struct")
			match := re.FindStringSubmatch(string(data))
			if len(match) > 0 {
				modelName := strings.TrimSpace(match[1])
				if modelName == "Base" {
					return nil
				}

				models = append(models, fmt.Sprintf("&models.%s{}", modelName))
			}
			return err
		}

		return nil
	})

	// in case someone removed this file, let's create it again
	if !files.FileExists("migrations.go") {
		err = ioutil.WriteFile("migrations.go", []byte(""), 0644)
		if err != nil {
			return err
		}
	}

	migrationsLocation, err := filepath.Abs("migrations.go")
	if err != nil {
		return err
	}

	contentString := MIGRATIONS_TEMPLATE
	contentString = strings.Replace(contentString, "{{Root}}", m.Config.Package, -1)
	contentString = strings.Replace(contentString, "{{Models}}", strings.Join(models, ", "), -1)

	content, err := format.Source([]byte(contentString))
	if err != nil {
		return err
	}

	err = m.saveFile(content, migrationsLocation)
	if err != nil {
		return err
	}

	return nil
}
