package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/alenn-m/rgen/v2/generator"
	"github.com/alenn-m/rgen/v2/generator/controller"
	"github.com/alenn-m/rgen/v2/generator/migration"
	"github.com/alenn-m/rgen/v2/generator/model"
	"github.com/alenn-m/rgen/v2/generator/parser"
	"github.com/alenn-m/rgen/v2/generator/repository"
	"github.com/alenn-m/rgen/v2/generator/service_init"
	"github.com/alenn-m/rgen/v2/generator/transport"
	"github.com/alenn-m/rgen/v2/util/config"
	helperLog "github.com/alenn-m/rgen/v2/util/log"
	"github.com/spf13/cobra"
)

var name string
var fields string
var actions string
var public bool
var onlyModel bool

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates API CRUD with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			log.Println(err.Error())
			return
		}

		conf, err := loadConfig(wd)
		if err != nil {
			log.Println(err.Error())
			return
		}

		// Initialize parser
		p := new(parser.Parser)

		p.OnlyModel = onlyModel
		p.Public = public

		err = p.Parse(name, fields, actions)
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = generate(p, conf)
		if err != nil {
			log.Println(err.Error())
			return
		}

		if p.OnlyModel {
			helperLog.Info(fmt.Sprintf("Model and migrations for '%s' are created", p.Name))
		} else {
			helperLog.Info(fmt.Sprintf("API for '%s' is created", p.Name))
		}
	},
}

func generate(p *parser.Parser, conf *config.Config) error {
	steps := []generator.GeneratorStep{&model.Model{}}

	if !p.OnlyModel {
		steps = append(steps,
			&repository.Repository{},
			&repository.RepositoryImplementation{},
			&controller.Controller{},
			&service_init.ServiceInit{},
			&transport.Transport{},
			&migration.Migration{},
		)
	}

	for _, step := range steps {
		err := step.Generate(p, conf)
		if err != nil {
			return err
		}

		err = step.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Resource name (required) --name='ModelName'")
	generateCmd.PersistentFlags().StringVarP(&fields, "fields", "f", "", "List of fields (required) --fields='Title:string, Description:string, UserID:int64'")
	generateCmd.PersistentFlags().StringVarP(&actions, "actions", "a", "", "CRUD actions --actions='index,create,show,update,delete'")
	generateCmd.PersistentFlags().BoolVar(&public, "public", false, "Public resource (default = false)")
	generateCmd.PersistentFlags().BoolVar(&onlyModel, "onlyModel", false, "Create only model (default = false)")

	_ = generateCmd.MarkFlagRequired("name")
	_ = generateCmd.MarkFlagRequired("fields")
}
