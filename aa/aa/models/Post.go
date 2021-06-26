package models

type PostID int64

type Post struct {
	ID     PostID `json:"id" db:"PostID"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"user_id"`
	User   *User  `json:"user"`

	Base
}
