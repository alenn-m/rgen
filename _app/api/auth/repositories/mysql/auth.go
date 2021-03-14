package mysql

import (
	"github.com/aa/aa/models"
	"github.com/jmoiron/sqlx"
)

type AuthDB struct {
	client *sqlx.DB
}

func NewAuthDB(client *sqlx.DB) *AuthDB {
	return &AuthDB{client: client}
}

func (a *AuthDB) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := a.client.Get(&user, "SELECT * FROM Users WHERE Email = ?", email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *AuthDB) UpdateToken(id models.UserID, token string) error {
	_, err := a.client.Exec("UPDATE Users SET ApiToken = ? WHERE UserID = ?", token, id)

	return err
}

func (a *AuthDB) FindByToken(token string) (*models.User, error) {
	var user models.User

	err := a.client.Get(&user, "SELECT * FROM Users WHERE ApiToken = ?", token)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *AuthDB) ClearToken(id models.UserID) error {
	_, err := a.client.Exec("UPDATE Users SET ApiToken = '' WHERE UserID = ?", id)

	return err
}

func (a *AuthDB) InsertUser(user *models.User) (models.UserID, error) {
	r, err := a.client.NamedExec("INSERT INTO Users (LastName, Password, ApiToken, Email, FirstName) VALUES (:LastName, :Password, :ApiToken, :Email, :FirstName)", user)
	if err != nil {
		return 0, err
	}

	insertID, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return models.UserID(insertID), nil
}
