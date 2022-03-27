package validators

import (
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(database *sqlx.DB) {
	db = database
}
