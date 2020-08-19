package test

type User struct {
	Base

	ID       int64  `json:"id" orm:"pk"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	ApiToken string `json:"api_token"`
}
