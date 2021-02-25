package user

import (
	"context"

	"github.com/aa/aa/models"
	authService "github.com/aa/aa/util/auth"
)

type DBRepository interface {
	ListAll(int) ([]models.User, int, error)
	FindByID(models.UserID) (*models.User, error)
	Insert(*models.User) (models.UserID, error)
	Update(*models.User) error
	Delete(models.UserID) error
}

type UsersController struct {
	db   DBRepository
	auth *authService.AuthService
}

type Repository interface {
	Index(context.Context, int) ([]models.User, int, error)
	Show(context.Context, models.UserID) (*models.User, error)
	Store(context.Context, *StoreReq) (models.UserID, error)
	Update(context.Context, *UpdateReq, models.UserID) error
	Delete(context.Context, models.UserID) error
}

func NewController(db DBRepository, auth *authService.AuthService) *UsersController {
	return &UsersController{
		db:   db,
		auth: auth,
	}
}
