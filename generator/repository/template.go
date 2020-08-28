package repository

const TEMPLATE = `
package {{Package}}

import (
	"context"

	"{{Root}}/models"
	authService "{{Root}}/util/auth"
)

type DBRepository interface {
	FindByID(int64) (*models.{{Model}}, error)
	ListAll() ([]models.{{Model}}, error)
	Insert(models.{{Model}}) (int64, error)
	Update(models.{{Model}}) error
	Delete(int64) error
}

type {{Controller}} struct {
	db   DBRepository
	auth *authService.AuthService
}

type Repository interface {
	Index(context.Context) ([]models.{{Model}}, error)
	Store(context.Context, *StoreReq) (int64, error)
	Show(context.Context, int64) (*models.{{Model}}, error)
	Update(context.Context, *UpdateReq, int64) error
	Delete(context.Context, int64) error
}

func NewController(db DBRepository, auth *authService.AuthService) *{{Controller}} {
	return &{{Controller}}{
		db:   db,
		auth: auth,
	}
}
`

const MYSQL_TEMPLATE = `
package mysql

import (
	"{{Root}}/models"
	"github.com/jinzhu/gorm"
)

type {{Model}}DB struct {
	client *gorm.DB
}

func New{{Model}}DB(client *gorm.DB) *{{Model}}DB {
	return &{{Model}}DB{client: client}
}

func (u *{{Model}}DB) FindByID(id int64) (*models.{{Model}}, error) {
	var item models.{{Model}}

	err := u.client.Where("id = ?", id).Find(&item).Error

	return &item, err
}

func (u *{{Model}}DB) ListAll() ([]models.{{Model}}, error) {
	var items []models.{{Model}}

	err := u.client.Find(&items).Error

	return items, err
}

func (u *{{Model}}DB) Insert(item models.{{Model}}) (int64, error) {
	err := u.client.Create(&item).Error

	return item.ID, err
}

func (u *{{Model}}DB) Update(item models.{{Model}}) error {
	return u.client.Model(&item).Updates(item).Error
}

func (u *{{Model}}DB) Delete(id int64) error {
	err := u.client.Where("id = ?", id).Delete(&models.{{Model}}{}).Error

	return err
}
`
