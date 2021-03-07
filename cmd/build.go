package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alenn-m/rgen/generator/migration"
	"github.com/alenn-m/rgen/generator/model"
	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/draft"
	"github.com/alenn-m/rgen/util/files"
	"github.com/alenn-m/rgen/util/log"
	"github.com/alenn-m/rgen/util/misc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// buildCmd represents the new command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds API from YAML file",
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			log.Error(err.Error())
			return
		}

		conf, err := loadConfig(wd)
		if err != nil {
			log.Error(err.Error())
			return
		}

		drft, err := loadDraft(wd)
		if err != nil {
			log.Error(err.Error())
			return
		}

		tables := []migration.PivotMigrationEntry{}

		for modelName, item := range drft.Models {
			fields := []parser.Field{}
			for k, v := range item.Properties {
				fields = append(fields, parser.Field{
					Key:   k,
					Value: v,
				})
			}

			if len(item.Actions) == 0 {
				item.Actions = misc.ACTIONS
			}

			p := parser.Parser{
				Name:          modelName,
				Fields:        fields,
				Relationships: item.Relationships,
				Actions:       actionsToUpper(item.Actions),
				OnlyModel:     item.OnlyModel,
				Public:        item.Public,
			}

			err := generate(&p, conf)
			if err != nil {
				log.Error(err.Error())
				return
			}

			for mn, relationship := range item.Relationships {
				if relationship == model.ManyToMany {
					tables = append(tables, migration.PivotMigrationEntry{
						TableOne: modelName,
						TableTwo: mn,
					})
				}
			}

		}

		// Generating pivot tables at the end to avoid duplicates
		// and to be sure referenced tables exists
		pm := new(migration.PivotMigration)
		pm.Init(tables)
		err = pm.Generate()
		if err != nil {
			log.Error(err.Error())
			return
		}

		log.Info(fmt.Sprintf("Built API for %d model(s)", len(drft.Models)))
	},
}

func actionsToUpper(input []string) []string {
	output := []string{}
	for _, item := range input {
		output = append(output, strings.ToUpper(item))
	}

	return output
}

func loadDraft(wd string) (*draft.Draft, error) {
	if !files.FileExists(fmt.Sprintf("%s/draft.yaml", wd)) {
		return nil, errors.New("draft.yaml not found")
	}

	draftData, err := ioutil.ReadFile(fmt.Sprintf("%s/draft.yaml", wd))
	if err != nil {
		return nil, err
	}

	var drft draft.Draft
	err = yaml.Unmarshal(draftData, &drft)
	if err != nil {
		return nil, err
	}

	return &drft, nil
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
