package models

type Post struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id" validate:"required"`
	Body     string `json:"body" validate:"required"`
}
