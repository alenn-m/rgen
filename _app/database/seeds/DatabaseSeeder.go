package seeds

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type Seeder interface {
	Execute() error
}

var seeders = []Seeder{
	// new(UsersSeeder),
}

type DatabaseSeeder struct {
	Client *sqlx.DB
}

func NewDatabaseSeeder(client *sqlx.DB) *DatabaseSeeder {
	return &DatabaseSeeder{Client: client}
}

func (d *DatabaseSeeder) Run() error {
	defer func() {
		if err := d.Client.Close(); err != nil {
			log.Fatalf("failed to close DB: %v\n", err)
		}
	}()

	for _, s := range seeders {
		if err := s.Execute(); err != nil {
			return err
		}
	}

	return nil
}
