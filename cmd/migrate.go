package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
	"github.com/spf13/cobra"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = "./database/migrations"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates database",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		if len(args) < 1 {
			flags.Usage()
			return
		}

		db, err := goose.OpenDBWithDriver("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_TABLE")),
		)
		if err != nil {
			log.Fatalf("goose: failed to open DB: %v\n", err)
		}

		defer func() {
			if err := db.Close(); err != nil {
				log.Fatalf("goose: failed to close DB: %v\n", err)
			}
		}()

		command := args[0]
		arguments := []string{}
		if len(args) > 3 {
			arguments = append(arguments, args[3:]...)
		}

		goose.SetSequential(true)
		if err := goose.Run(command, db, dir, arguments...); err != nil {
			log.Fatalf("goose %v: %v", command, err)
		}
	},
}
