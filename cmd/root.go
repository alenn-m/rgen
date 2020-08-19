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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var rootDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&rootDir, "dir", "", "root directory of the project (default is current working directory)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".generate" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".generate")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
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
