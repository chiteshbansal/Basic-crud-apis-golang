// The repository package provides functionalities for data manipulation and validation on the User model.
package repository

import (
	db "first-api/internal/database"
	model "first-api/internal/models"

	"github.com/jinzhu/gorm"
)

// UserStorer is the interface that wraps the methods for manipulating and validating User data.
type PostStorer interface {
	CreatePost(post *model.Post) error
	GetAllPosts(post *[]model.Post) error
	DeletePost(post *model.Post, query string) error
	GetPost(post *model.Post, query string) error
	UpdatePost(post *model.Post) error
}

// UserStore is a concrete implementation of the UserStorer interface.
type PostStore struct{}

// CreateUser creates a user in the database.
func (ps *PostStore) CreatePost(post *model.Post) error {
	if err := db.DB.Create(post).Error; err != nil {
		return &CustomError{Message: "Post cannot be Created "}
	}
	return nil
}

// GetAllUsers fetches all user data from the database.
func (ps *PostStore) GetAllPosts(post *[]model.Post) error {
	if err := db.DB.Preload("Creator",selectRequiredFields).Find(post).Error; err != nil {
		return &CustomError{Message: "Posts cannot be fetched  "}
	}
	return nil
}

// DeletePPost deletes a post by the given ID from the database.
func (ps *PostStore) DeletePost(post *model.Post, id string) (err error) {
	if err = db.DB.Where("id = ?", id).Delete(post).Error; err != nil {
		return &CustomError{Message: "Delete Post Failed "}
	}
	return nil
}

// getPost fetches a post by the given query from the database.
func (ps *PostStore) GetPost(post *model.Post, query string) (err error) {
	if err = db.DB.Where(query).Preload("Creator", selectRequiredFields).Preload("Comments.Author", selectRequiredFields).Find(post).Error; err != nil {
		return &CustomError{Message: "Post Not found "}
	}
	return nil
}

func (ps *PostStore) UpdatePost(post *model.Post) (err error) {
	if err = db.DB.Save(post).Error; err != nil {
		return &CustomError{Message: "Post Update Failed ! Try again"}
	}
	return nil
}

func selectRequiredFields(db *gorm.DB) *gorm.DB {
	return db.Select("id, name")
}
