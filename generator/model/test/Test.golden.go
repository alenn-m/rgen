package models

type Test struct {
	Base

	ID   int64  `json:"id" orm:"pk"`
	Name string `json:"name"`
}
