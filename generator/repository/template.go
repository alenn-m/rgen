package repository

const DBR_SHOW = "FindByID(models.{{Model}}ID) (*models.{{Model}}, error)"
const DBR_INDEX = "ListAll(int) ([]models.{{Model}}, int, error)"
const DBR_CREATE = "Insert(*models.{{Model}}) (models.{{Model}}ID, error)"
const DBR_UPDATE = "Update(*models.{{Model}}) error"
const DBR_DELETE = "Delete(models.{{Model}}ID) error"

const R_SHOW = "Show(context.Context, models.{{Model}}ID) (*models.{{Model}}, error)"
const R_INDEX = "Index(context.Context, int) ([]models.{{Model}}, int, error)"
const R_CREATE = "Store(context.Context, *StoreReq) (models.{{Model}}ID, error)"
const R_UPDATE = "Update(context.Context, *UpdateReq, models.{{Model}}ID) error"
const R_DELETE = "Delete(context.Context, models.{{Model}}ID) error"

const TEMPLATE_AUTH = `
package {{Package}}

import (
	"context"

	"{{Root}}/models"
	authService "{{Root}}/util/auth"
)

type DBRepository interface {
    {{DBRepositoryActions}}
}

type {{Controller}} struct {
	db   DBRepository
	auth *authService.AuthService
}

type Repository interface {
    {{RepositoryActions}}
}

func NewController(db DBRepository, auth *authService.AuthService) *{{Controller}} {
	return &{{Controller}}{
		db:   db,
		auth: auth,
	}
}
`

const TEMPLATE_NO_AUTH = `
package {{Package}}

import (
	"context"

	"{{Root}}/models"
)

type DBRepository interface {
    {{DBRepositoryActions}}
}

type {{Controller}} struct {
	db DBRepository
}

type Repository interface {
    {{RepositoryActions}}
}

func NewController(db DBRepository) *{{Controller}} {
	return &{{Controller}}{
		db: db,
	}
}
`

const MYSQL_TEMPLATE_HEADER = `
package mysql

import (
	"{{Root}}/models"
    "{{Root}}/util/paginate"
	"github.com/jinzhu/gorm"
)

type {{Model}}DB struct {
	client *gorm.DB
}

func New{{Model}}DB(client *gorm.DB) *{{Model}}DB {
	return &{{Model}}DB{client: client}
}`

const MYSQL_TEMPLATE_SHOW = `func (u *{{Model}}DB) FindByID(id models.{{Model}}ID) (*models.{{Model}}, error) {
	var item models.{{Model}}

	err := u.client.Where("id = ?", id).Find(&item).Error

	return &item, err
}`

const MYSQL_TEMPLATE_INDEX = `func (u *{{Model}}DB) ListAll(page int) ([]models.{{Model}}, int, error) {
	var items []models.{{Model}}

    err := paginate.Paginate(u.client.New(), page).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

    count := 0
	err = u.client.Model(&models.{{Model}}{}).Count(&count).Error

	return items, count, err
}`

const MYSQL_TEMPLATE_CREATE = `func (u *{{Model}}DB) Insert(item *models.{{Model}}) (models.{{Model}}ID, error) {
	err := u.client.Create(&item).Error

	return item.ID, err
}`

const MYSQL_TEMPLATE_UPDATE = `func (u *{{Model}}DB) Update(item *models.{{Model}}) error {
	return u.client.Model(&item).Updates(item).Error
}`

const MYSQL_TEMPLATE_DELETE = `func (u *{{Model}}DB) Delete(id models.{{Model}}ID) error {
	return u.client.Where("id = ?", id).Delete(&models.{{Model}}{}).Error
}`
