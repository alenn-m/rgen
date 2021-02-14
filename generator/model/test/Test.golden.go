package test

type Test struct {
	ID   int64  `json:"id" orm:"pk"`
	Name string `json:"name"`
	User User   `json:"user"`

	Base
}
