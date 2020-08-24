package models

type Test struct {
	ID   int64  `json:"id" orm:"pk"`
	Name string `json:"name"`

	Base
}
