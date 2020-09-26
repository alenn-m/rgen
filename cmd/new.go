package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/draft"
	"github.com/alenn-m/rgen/util/files"
	"github.com/alenn-m/rgen/util/log"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Initializes the REST API",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Creating new application")

		wd, err := os.Getwd()
		if err != nil {
			log.Error(err.Error())
			return
		}

		vscItems := []string{"github.com", "bitbucket.org"}

		promptVSC := promptui.Select{
			Label: "Select VSC",
			Items: vscItems,
		}

		_, vsc, err := promptVSC.Run()
		if err != nil {
			log.Error(err.Error())
			return
		}

		pDomain := promptui.Prompt{
			Label: "Domain",
		}

		domain, err := pDomain.Run()
		if err != nil {
			log.Error(err.Error())
			return
		}

		pPackage := promptui.Prompt{
			Label: "Package",
		}

		pkg, err := pPackage.Run()
		if err != nil {
			log.Error(err.Error())
			return
		}

		rootPackage := fmt.Sprintf("%s/%s/%s", vsc, domain, pkg)

		destination := fmt.Sprintf("%s/%s", wd, pkg)

		err = files.CopyDir(fmt.Sprintf("%s/src/github.com/alenn-m/rgen/_app", os.Getenv("GOPATH")), destination)
		if err != nil {
			log.Error(err.Error())
			return
		}

		config := config.Config{
			Package:        rootPackage,
			AutoMigrations: true,
		}

		configData, err := yaml.Marshal(config)
		if err != nil {
			log.Error(err.Error())
			return
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/config.yaml", destination), configData, 0644)
		if err != nil {
			log.Error(err.Error())
			return
		}

		drft := draft.Draft{
			Models: map[string]draft.ModelOptions{
				"User": {
					Properties: map[string]string{
						"FirstName": "string",
						"LastName":  "string",
						"Email":     "string",
						"ApiToken":  "string",
						"Password":  "string",
					},
					Validation: map[string][]string{
						"FirstName": {"Required"},
						"LastName":  {"Required"},
						"Email":     {"Required", "Email"},
					},
					Actions:       []string{ACTION_ALL},
					Relationships: nil,
				},
			},
		}

		draftData, err := yaml.Marshal(drft)
		if err != nil {
			log.Error(err.Error())
			return
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/draft.yaml", destination), draftData, 0644)
		if err != nil {
			log.Error(err.Error())
			return
		}

		err = os.Rename(fmt.Sprintf("%s/.env.example", destination), fmt.Sprintf("%s/.env", destination))
		if err != nil {
			log.Error(err.Error())
			return
		}

		err = filepath.Walk(destination, func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				data = []byte(strings.Replace(string(data), "{{Root}}", rootPackage, -1))

				err = ioutil.WriteFile(path, data, 0644)

				return err
			}

			return err
		})
		if err != nil {
			fmt.Println("sdfsfsd")
			log.Error(err.Error())
			return
		}

		log.Notice("Application created!")
		log.Info("Add your models definitions to draft.yaml file!")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
