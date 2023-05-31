package repository

import (
	"context"
	model "first-api/api/Models"
	db "first-api/database"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserStorer interface {
	CreateUser(ctx context.Context, user *model.User) error
	Validate(user model.User) error
}
type UserStore struct{}

type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func (u *UserStore) CreateUser(ctx context.Context, user *model.User) error {
	if err := db.DB.Create(user).Error; err != nil {
		return &CustomError{Message: "User cannot be Created "}
	}
	return nil
}

func (us *UserStore) Validate(u model.User) error {
	return validation.ValidateStruct(&u, validation.Field(&u.Name, validation.Required, validation.Length(5, 20)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Phone, validation.Required, validation.Length(10, 10)),
		validation.Field(&u.Address, validation.Required, validation.Length(10, 50)),
	)
}

// // get all users Fetch all user data
func (us *UserStore) GetAllUsers(user *[]model.User) error {
	if err := db.DB.Find(user).Error; err != nil {
		return &CustomError{Message: "Users cannot be fetched  "}
	}
	return nil
}

// getuserById

func (us *UserStore) GetUser(user *model.User, query string) (err error) {
	if err = db.DB.Where(query).Find(user).Error; err != nil {
		return &CustomError{Message: "User Not found "}
	}
	return nil
}

// update user

func (us *UserStore) UpdateUser(user *model.User, id string) (err error) {
	if err = db.DB.Save(user).Error; err != nil {
		return &CustomError{Message: "User Update Failed! Try Again"}
	}
	return nil
}

// Delete User

func (us *UserStore) DeleteUser(user *model.User, id string) (err error) {
	if err = db.DB.Where("id = ?", id).Delete(user).Error; err != nil {
		return &CustomError{Message: "Delete User Failed "}
	}
	return nil
}
