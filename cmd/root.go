package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alenn-m/rgen/util/config"
	"github.com/alenn-m/rgen/util/draft"
	"github.com/alenn-m/rgen/util/files"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
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

func init() {
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	// Find home directory.
// 	home, err := homedir.Dir()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
//
// 	// Search config in home directory with name ".generate" (without extension).
// 	viper.AddConfigPath(home)
// 	viper.SetConfigName(".generate")
//
// 	viper.AutomaticEnv() // read in environment variables that match
//
// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }

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
