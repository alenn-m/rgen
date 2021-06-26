package user

import (
	"context"

	"github.com/aa/aa/models"
)

func (u *UsersController) Index(c context.Context, page int) ([]models.User, int, error) {
	items, total, err := u.db.ListAll(page)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (u *UsersController) Store(c context.Context, r *StoreReq) (models.UserID, error) {
	id, err := u.db.Insert(&models.User{
		ApiToken:  r.ApiToken,
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Password:  r.Password,
	})

	return id, err
}

func (u *UsersController) Show(c context.Context, id models.UserID) (*models.User, error) {
	return u.db.FindByID(id)
}

func (u *UsersController) Update(c context.Context, r *UpdateReq, id models.UserID) error {
	return u.db.Update(&models.User{
		ID:        id,
		ApiToken:  r.ApiToken,
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Password:  r.Password,
	})
}

func (u *UsersController) Delete(c context.Context, id models.UserID) error {
	return u.db.Delete(id)
}
