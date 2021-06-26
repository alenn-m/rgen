package post

import (
	"context"

	"github.com/aa/aa/models"
	authService "github.com/aa/aa/util/auth"
)

type DBRepository interface {
	ListAll(int) ([]models.Post, int, error)
	FindByID(models.PostID) (*models.Post, error)
	Insert(*models.Post) (models.PostID, error)
	Update(*models.Post) error
	Delete(models.PostID) error
}

type PostsController struct {
	db   DBRepository
	auth *authService.AuthService
}

type Repository interface {
	Index(context.Context, int) ([]models.Post, int, error)
	Show(context.Context, models.PostID) (*models.Post, error)
	Store(context.Context, *StoreReq) (models.PostID, error)
	Update(context.Context, *UpdateReq, models.PostID) error
	Delete(context.Context, models.PostID) error
}

func NewController(db DBRepository, auth *authService.AuthService) *PostsController {
	return &PostsController{
		db:   db,
		auth: auth,
	}
}
