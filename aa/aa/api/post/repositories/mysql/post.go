package mysql

import (
	"github.com/aa/aa/models"
	"github.com/aa/aa/util/paginate"
	"github.com/jmoiron/sqlx"
)

type PostDB struct {
	client *sqlx.DB
}

func NewPostDB(client *sqlx.DB) *PostDB {
	return &PostDB{client: client}
}

func (u *PostDB) FindByID(id models.PostID) (*models.Post, error) {
	var item models.Post

	err := u.client.Get(&item, "SELECT * FROM Posts WHERE PostID = ?", id)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (u *PostDB) ListAll(page int) ([]models.Post, int, error) {
	var items []models.Post

	query := paginate.Paginate("SELECT * FROM Posts", page)
	err := u.client.Select(&items, query)
	if err != nil {
		return nil, 0, err
	}

	count := 0
	err = u.client.Get(&count, "SELECT COUNT(PostID) FROM Posts")

	return items, count, err
}

func (u *PostDB) Insert(item *models.Post) (models.PostID, error) {
	r, err := u.client.NamedExec("INSERT INTO Posts (Title, Body, UserID) VALUES (:Title, :Body, :UserID)", item)
	if err != nil {
		return 0, err
	}

	insertID, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return models.PostID(insertID), nil
}

func (u *PostDB) Update(item *models.Post) error {
	_, err := u.client.NamedExec(`UPDATE Posts
        SET Title = :Title, Body = :Body, UserID = :UserID
        WHERE PostID = :PostID
    `, item)

	return err
}

func (u *PostDB) Delete(id models.PostID) error {
	_, err := u.client.Exec("DELETE FROM Posts WHERE PostID = ?", id)

	return err
}
