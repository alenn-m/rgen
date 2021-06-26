package post

import (
	"context"

	"github.com/aa/aa/models"
)

func (u *PostsController) Index(c context.Context, page int) ([]models.Post, int, error) {
	items, total, err := u.db.ListAll(page)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (u *PostsController) Store(c context.Context, r *StoreReq) (models.PostID, error) {
	id, err := u.db.Insert(&models.Post{
		Title:  r.Title,
		Body:   r.Body,
		UserID: r.UserID,
	})

	return id, err
}

func (u *PostsController) Show(c context.Context, id models.PostID) (*models.Post, error) {
	return u.db.FindByID(id)
}

func (u *PostsController) Update(c context.Context, r *UpdateReq, id models.PostID) error {
	return u.db.Update(&models.Post{
		ID:     id,
		Title:  r.Title,
		Body:   r.Body,
		UserID: r.UserID,
	})
}

func (u *PostsController) Delete(c context.Context, id models.PostID) error {
	return u.db.Delete(id)
}
