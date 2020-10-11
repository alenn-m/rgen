package cmd

import (
	"log"
	"os"

	"github.com/alenn-m/rgen/generator/controller"
	"github.com/alenn-m/rgen/generator/model"
	"github.com/alenn-m/rgen/generator/parser"
	"github.com/alenn-m/rgen/generator/repository"
	"github.com/alenn-m/rgen/generator/service_init"
	"github.com/alenn-m/rgen/generator/transport"
	"github.com/alenn-m/rgen/util/config"
	"github.com/spf13/cobra"
)

var name string
var fields string
var actions string

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
	},
}

func generate(p *parser.Parser, conf *config.Config) error {
	// Generate model
	m := new(model.Model)
	m.Init(&model.Input{
		Name:          p.Name,
		Fields:        p.Fields,
		Relationships: p.Relationships,
	}, conf)
	err := m.Generate()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if !p.SkipController {
		// Generate repositories
		r := new(repository.Repository)
		r.Init(&repository.Input{
			Name:    p.Name,
			Actions: p.Actions,
			Public:  p.Public,
		}, conf)
		err = r.Generate()
		if err != nil {
			log.Println(err.Error())
			return err
		}

		// Generate controller
		c := new(controller.Controller)
		c.Init(&controller.Input{
			Name:    p.Name,
			Fields:  p.Fields,
			Actions: p.Actions,
		}, conf)
		err = c.Generate()
		if err != nil {
			log.Println(err.Error())
			return err
		}

		// Generate services
		serviceInit := new(service_init.ServiceInit)
		serviceInit.Init(&service_init.Input{
			Name:   p.Name,
			Public: p.Public,
		}, conf)
		err = serviceInit.Generate()
		if err != nil {
			return err
		}
	}

	// Generate transport layer
	t := new(transport.Transport)
	t.Init(&transport.Input{
		Name:    p.Name,
		Fields:  p.Fields,
		Actions: p.Actions,
	}, conf)
	err = t.Generate()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Initiate auto migrations
	if conf.AutoMigrations {
		err = m.SetupAutoMigration()
		if err != nil {
			log.Println(err.Error())
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

	_ = generateCmd.MarkFlagRequired("name")
	_ = generateCmd.MarkFlagRequired("fields")
}
