package entities

import "time"

// Comment represents the domain entity for a blog comment
type Comment struct {
	ID         string    `json:"id"`
	PostID     string    `json:"post_id"`
	AuthorID   string    `json:"author_id"`
	Content    string    `json:"content"`
	IsApproved bool      `json:"is_approved"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
