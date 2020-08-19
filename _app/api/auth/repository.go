package auth

import (
	"context"

	"{{Root}}/models"
	authService "{{Root}}/util/auth"
)

type Repository interface {
	Login(context.Context, string, string) (*models.User, error)
	Signup(context.Context, *SignupReq) (int64, error)
	Logout(context.Context) error
}

type DBRepository interface {
	InsertUser(*models.User) (int64, error)
	FindByEmailAndPassword(string, string) (*models.User, error)
	FindByToken(string) (*models.User, error)
	UpdateToken(int64, string) error
	ClearToken(int64) error
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
