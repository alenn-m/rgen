package models

type Test struct {
	ID    int64  `json:"id" orm:"pk"`
	Name  string `json:"name"`
	User  User   `json:"user"`
	Posts []Post `json:"posts" gorm:"many2many:post_tests"`

	Base
}
