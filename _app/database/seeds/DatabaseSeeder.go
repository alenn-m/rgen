package seeds

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type Seeder interface {
	Execute(*sqlx.DB) error
}

var seeders = []Seeder{
	new(UserSeeder),
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
		if err := s.Execute(d.Client); err != nil {
			return err
		}
	}

	return nil
}
