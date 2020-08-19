package auth

import (
	"context"
	"errors"

	"{{Root}}/models"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

func (u *AuthController) Login(c context.Context, email, password string) (*models.User, error) {
	user, err := u.db.FindByEmailAndPassword(email, password)
	if err != nil {
		return nil, ErrInvalidCredentials
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

func (u *AuthController) Signup(c context.Context, r *SignupReq) (int64, error) {
	user := models.User{
		Email:    r.Email,
		Password: r.Password,
	}

	token, err := u.auth.GenerateToken(&user)
	if err != nil {
		return 0, err
	}
	user.ApiToken = token.Raw

	return u.db.InsertUser(&user)
}
