package validator

import (
	model "first-api/internal/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

// CreatePost represents the data structure for creating a post.
type CreatePost struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Creator     model.User      `json:"creator" gorm:"foreignKey:CreatorId"`
	Comments    []model.Comment `json:"comments"`
}

// Validate performs the validation for the CreatePost struct.
// It checks the following rules:
// - Title is required and must have a length between 5 and 100 characters.
// - Description is required and must have a length between 10 and 500 characters.
func (c CreatePost) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Title, validation.Required, validation.Length(5, 100)),
		validation.Field(&c.Description, validation.Required, validation.Length(10, 500)),
	)
}

// UpdatePost represents the data structure for updating a post.
type UpdatePost struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Creator     model.User      `json:"creator" gorm:"foreignKey:CreatorId"`
	Comments    []model.Comment `json:"comments"`
	CreatorId   uint            `json:"creator_id"`
}

// Validate performs the validation for the UpdatePost struct.
// It checks the following rules:
// - Title is required and must have a length between 5 and 100 characters.
// - Description is required and must have a length between 10 and 500 characters.
// - CreatorId is required.
func (c UpdatePost) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Title, validation.Required, validation.Length(5, 100)),
		validation.Field(&c.Description, validation.Required, validation.Length(10, 500)),
		validation.Field(&c.CreatorId, validation.Required),
	)
}

// Comment represents the data structure for a comment.
type Comment struct {
	PostID   string `json:"post_id"`
	AuthorID uint   `json:"author_id"`
	Content  string `json:"content"`
}

// Validate performs the validation for the Comment struct.
// It checks the following rules:
// - PostID is required.
// - AuthorID is required and must be greater than 0.
// - Content is required and must not be empty.
func (c Comment) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.PostID, validation.Required),
		validation.Field(&c.AuthorID, validation.Required, validation.Min(uint(1))),
		validation.Field(&c.Content, validation.Required, validation.Length(1, 0)),
	)
}
