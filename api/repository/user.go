// The repository package provides functionalities for data manipulation and validation on the User model.
package repository

import (
	"context"
	model "first-api/api/Models"
	db "first-api/database"
)

// UserStorer is the interface that wraps the methods for manipulating and validating User data.
type UserStorer interface {
	CreateUser(ctx context.Context, user *model.User) error
	Validate(user model.User) error
}

// UserStore is a concrete implementation of the UserStorer interface.
type UserStore struct{}

// CustomError represents a custom error with a message.
type CustomError struct {
	Message string
}

// Error returns the error message.
func (e *CustomError) Error() string {
	return e.Message
}

// CreateUser creates a user in the database.
func (u *UserStore) CreateUser(ctx context.Context, user *model.User) error {
	if err := db.DB.Create(user).Error; err != nil {
		return &CustomError{Message: "User cannot be Created "}
	}
	return nil
}


// GetAllUsers fetches all user data from the database.
func (us *UserStore) GetAllUsers(user *[]model.User) error {
	if err := db.DB.Find(user).Error; err != nil {
		return &CustomError{Message: "Users cannot be fetched  "}
	}
	return nil
}

// GetUser fetches a user by the given query from the database.
func (us *UserStore) GetUser(user *model.User, query string) (err error) {
	if err = db.DB.Where(query).Find(user).Error; err != nil {
		return &CustomError{Message: "User Not found "}
	}
	return nil
}

// UpdateUser updates a user in the database.
func (us *UserStore) UpdateUser(user *model.User, id string) (err error) {
	if err = db.DB.Save(user).Error; err != nil {
		return &CustomError{Message: "User Update Failed! Try Again"}
	}
	return nil
}

// DeleteUser deletes a user by the given ID from the database.
func (us *UserStore) DeleteUser(user *model.User, id string) (err error) {
	if err = db.DB.Where("id = ?", id).Delete(user).Error; err != nil {
		return &CustomError{Message: "Delete User Failed "}
	}
	return nil
}
