package auth

import (
	"context"

	"github.com/alenn-m/myApp/models"
	"github.com/dgrijalva/jwt-go"
)

type AuthInterface interface {
	GenerateToken(*models.User) (*jwt.Token, error)
	GetLoggedInUser(context.Context) *models.User
	GetAuthUser(string) (*models.User, error)
	SetAuthUser(string, *models.User) error
	Logout(string) error
}
