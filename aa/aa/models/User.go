package models

type UserID int64

type User struct {
	ID        UserID `json:"id" db:"UserID"`
	ApiToken  string `json:"api_token"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`

	Base
}
