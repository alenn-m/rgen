package auth

import (
	"context"
	"errors"
	"time"

	"github.com/alenn-m/myApp/models"
	"github.com/alenn-m/myApp/util/cache"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

const (
	JwtSecret = "Eo3Xc5fBML"
	AuthUser  = "AUTH_USER"
)

var TokenAuth = jwtauth.New("HS256", []byte(JwtSecret), nil)

type AuthService struct {
	cache cache.CacheService
}

func NewAuthService(cache cache.CacheService) *AuthService {
	return &AuthService{
		cache: cache,
	}
}

func (a AuthService) GenerateToken(user *models.User) (*jwt.Token, error) {
	type MyClaims struct {
		ID    int64
		Email string
		Name  string
		jwt.StandardClaims
	}

	token, _, err := TokenAuth.Encode(
		MyClaims{
			ID:    user.ID,
			Email: user.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		},
	)

	return token, err
}

func (a AuthService) GetLoggedInUser(ctx context.Context) *models.User {
	user := ctx.Value(AuthUser)
	if user != nil {
		return user.(*models.User)
	}

	return nil
}

func (a AuthService) GetAuthUser(key string) (*models.User, error) {
	data, ok := a.cache.Get(key)
	if !ok {
		return nil, errors.New("user not found")
	}

	user := data.(*models.User)

	return user, nil
}

func (a AuthService) SetAuthUser(key string, user *models.User) error {
	a.cache.Set(key, user)

	return nil
}

func (a AuthService) Logout(key string) error {
	a.cache.Delete(key)

	return nil
}
