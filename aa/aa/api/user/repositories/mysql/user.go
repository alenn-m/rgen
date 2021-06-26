package mysql

import (
	"github.com/aa/aa/models"
	"github.com/aa/aa/util/paginate"
	"github.com/jmoiron/sqlx"
)

type UserDB struct {
	client *sqlx.DB
}

func NewUserDB(client *sqlx.DB) *UserDB {
	return &UserDB{client: client}
}

func (u *UserDB) FindByID(id models.UserID) (*models.User, error) {
	var item models.User

	err := u.client.Get(&item, "SELECT * FROM Users WHERE UserID = ?", id)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (u *UserDB) ListAll(page int) ([]models.User, int, error) {
	var items []models.User

	query := paginate.Paginate("SELECT * FROM Users", page)
	err := u.client.Select(&items, query)
	if err != nil {
		return nil, 0, err
	}

	count := 0
	err = u.client.Get(&count, "SELECT COUNT(UserID) FROM Users")

	return items, count, err
}

func (u *UserDB) Insert(item *models.User) (models.UserID, error) {
	r, err := u.client.NamedExec("INSERT INTO Users (ApiToken, Email, FirstName, LastName, Password) VALUES (:ApiToken, :Email, :FirstName, :LastName, :Password)", item)
	if err != nil {
		return 0, err
	}

	insertID, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return models.UserID(insertID), nil
}

func (u *UserDB) Update(item *models.User) error {
	_, err := u.client.NamedExec(`UPDATE Users
        SET ApiToken = :ApiToken, Email = :Email, FirstName = :FirstName, LastName = :LastName, Password = :Password
        WHERE UserID = :UserID
    `, item)

	return err
}

func (u *UserDB) Delete(id models.UserID) error {
	_, err := u.client.Exec("DELETE FROM Users WHERE UserID = ?", id)

	return err
}
