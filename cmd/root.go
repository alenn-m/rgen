package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alenn-m/rgen/v2/util/config"
	"github.com/alenn-m/rgen/v2/util/files"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "RGEN - Go (GoLang) REST code generator",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfig(wd string) (*config.Config, error) {
	if !files.FileExists(fmt.Sprintf("%s/config.yaml", wd)) {
		return nil, errors.New("config.yaml not found")
	}

	configData, err := ioutil.ReadFile(fmt.Sprintf("%s/config.yaml", wd))
	if err != nil {
		return nil, err
	}

	var conf config.Config
	err = yaml.Unmarshal(configData, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
