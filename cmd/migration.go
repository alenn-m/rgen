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
	sequential    bool
	migrationType string
	flags         = flag.NewFlagSet("goose", flag.ExitOnError)
	dir           = "./database/migrations"
)

func init() {
	rootCmd.AddCommand(migrationCmd)

	migrationCmd.PersistentFlags().BoolVar(&sequential, "sequential", false, "Create migration in sequential order")
	migrationCmd.PersistentFlags().StringVarP(&migrationType, "migration-type", "t", "sql", "Migration type (sql|go)")
}

// migrationCmd represents the migrate command
var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Manages database migrations",
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
		if len(args) > 1 {
			arguments = append(arguments, args[1:]...)
			arguments = append(arguments, migrationType)
		}

		goose.SetSequential(sequential)
		if err := goose.Run(command, db, dir, arguments...); err != nil {
			log.Fatalf("goose %v: %v", command, err)
		}
	},
}
