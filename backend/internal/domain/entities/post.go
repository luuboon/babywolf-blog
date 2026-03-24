package entities

import "time"

// Post represents the domain entity for a blog post
type Post struct {
	ID            string    `json:"id"`
	AuthorID      string    `json:"author_id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Content       string    `json:"content"`
	CoverImageURL string    `json:"cover_image_url,omitempty"`
	Published     bool      `json:"published"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
