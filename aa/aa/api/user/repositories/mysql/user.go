package mysql

import (
	"github.com/aa/aa/models"
	"github.com/aa/aa/util/paginate"
	"github.com/jinzhu/gorm"
)

type UserDB struct {
	client *gorm.DB
}

func NewUserDB(client *gorm.DB) *UserDB {
	return &UserDB{client: client}
}

func (u *UserDB) FindByID(id models.UserID) (*models.User, error) {
	var item models.User

	err := u.client.Where("id = ?", id).Find(&item).Error

	return &item, err
}

func (u *UserDB) ListAll(page int) ([]models.User, int, error) {
	var items []models.User

	err := paginate.Paginate(u.client.New(), page).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	count := 0
	err = u.client.Model(&models.User{}).Count(&count).Error

	return items, count, err
}

func (u *UserDB) Insert(item *models.User) (models.UserID, error) {
	err := u.client.Create(&item).Error

	return item.ID, err
}

func (u *UserDB) Update(item *models.User) error {
	return u.client.Model(&item).Updates(item).Error
}

func (u *UserDB) Delete(id models.UserID) error {
	return u.client.Where("id = ?", id).Delete(&models.User{}).Error
}
