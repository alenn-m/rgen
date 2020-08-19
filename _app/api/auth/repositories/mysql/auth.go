package mysql

import (
	"{{Root}}/models"
	"github.com/jinzhu/gorm"
)

type AuthDB struct {
	client *gorm.DB
}

func NewAuthDB(client *gorm.DB) *AuthDB {
	return &AuthDB{client: client}
}

func (a *AuthDB) FindByEmailAndPassword(email, password string) (*models.User, error) {
	var user models.User

	err := a.client.Where("email = ?", email).Where("password = ?", password).Find(&user).Error

	return &user, err
}

func (a *AuthDB) UpdateToken(id int64, token string) error {
	err := a.client.Model(models.User{}).Where("id = ?", id).Update("api_token", token).Error

	return err
}

func (a *AuthDB) FindByToken(token string) (*models.User, error) {
	var user models.User

	err := a.client.Where("api_token = ?", token).Find(&user).Error

	return &user, err
}

func (a *AuthDB) ClearToken(id int64) error {
	err := a.client.Model(models.User{}).Where("id = ?", id).Update("api_token", "").Error

	return err
}

func (a *AuthDB) InsertUser(user *models.User) (int64, error) {
	err := a.client.Save(user).Error
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
