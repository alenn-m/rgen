package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/test/testApp/models"
	"github.com/test/testApp/util/paginate"
)

type TestDB struct {
	client *sqlx.DB
}

func NewTestDB(client *sqlx.DB) *TestDB {
	return &TestDB{client: client}
}

func (u *TestDB) FindByID(id models.TestID) (*models.Test, error) {
	var item models.Test

	err := u.client.Get(&item, "SELECT * FROM Tests WHERE TestID = ?", id)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (u *TestDB) ListAll(page int) ([]models.Test, int, error) {
	var items []models.Test

	query := paginate.Paginate("SELECT * FROM Tests", page)
	err := u.client.Select(&items, query)
	if err != nil {
		return nil, 0, err
	}

	count := 0
	err = u.client.Get(&count, "SELECT COUNT(TestID) FROM Tests")

	return items, count, err
}

func (u *TestDB) Insert(item *models.Test) (models.TestID, error) {
	r, err := u.client.NamedExec("INSERT INTO Tests () VALUES ()", item)
	if err != nil {
		return 0, err
	}

	insertID, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return models.TestID(insertID), nil
}

func (u *TestDB) Update(item *models.Test) error {
	_, err := u.client.NamedExec(`UPDATE Tests
        SET 
        WHERE TestID = :TestID
    `, item)

	return err
}

func (u *TestDB) Delete(id models.TestID) error {
	_, err := u.client.Exec("DELETE FROM Tests WHERE TestID = ?", id)

	return err
}
