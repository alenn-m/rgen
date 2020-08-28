package auth

import (
    "context"
	"errors"

	"{{Root}}/models"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

func (u *AuthController) Login(c context.Context, email, password string) (*models.User, error) {
	user, err := u.db.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	token, err := u.auth.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	user.ApiToken = token.Raw
	err = u.db.UpdateToken(user.ID, token.Raw)
	if err != nil {
		return nil, err
	}

	// Caching the user
	err = u.auth.SetAuthUser(token.Raw, user)

	return user, err
}

func (u *AuthController) Logout(c context.Context) error {
	user := u.auth.GetLoggedInUser(c)
	if user != nil {
		u.auth.Logout(user.ApiToken)

		return u.db.ClearToken(user.ID)
	}

	return nil
}
