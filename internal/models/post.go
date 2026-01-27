package models

import (
	"time"

	"github.com/google/uuid"
)

var (
	postDB PostDB
)

func SetPostDB(d PostDB) {
	postDB = d
}

type PostDB interface {
	SavePost(*Post) error
	UpdatePost(*Post) error
	DeletePost(*Post) error
	GetPost(*Post) error
	ListPosts() ([]Post, error)
}

type Post struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      `json:"user"`
}

func NewPost(c string, u User) *Post {
	return &Post{
		Content: c,
		User:    u,
	}
}

func (p *Post) Save() error {
	return postDB.SavePost(p)
}

func (p *Post) Update() error {
	return nil
}

func (p *Post) Delete() error {
	return postDB.DeletePost(p)
}

func ListPost() ([]Post, error) {
	return postDB.ListPosts()
}
