package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/util/log"
	"github.com/alenn-m/rgen/util/misc"
	"github.com/spf13/cobra"
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

func init() {
	rootCmd.AddCommand(buildCmd)
}
