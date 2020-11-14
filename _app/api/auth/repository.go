package auth

import (
    "context"

	"{{Root}}/models"
	authService "{{Root}}/util/auth"
)

type Repository interface {
	Login(context.Context, string, string) (*models.User, error)
	Logout(context.Context) error
}

type DBRepository interface {
	InsertUser(*models.User) (models.UserID, error)
	FindByEmail(string) (*models.User, error)
	FindByToken(string) (*models.User, error)
	UpdateToken(models.UserID, string) error
	ClearToken(models.UserID) error
}

type AuthController struct {
	db   DBRepository
	auth *authService.AuthService
}

func NewController(db DBRepository, auth *authService.AuthService) *AuthController {
	return &AuthController{
		db:   db,
		auth: auth,
	}
}
