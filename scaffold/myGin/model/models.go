package model

type Info struct {
	Id       int64
	Created  int64
	Title    string
	Tags     string `json:"-,omitempty"`
	Category string
	NewTags  []string `db:"-"`
}
