// Package validator provides validation rules for various structs.
package validator

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// CreateUser represents the data structure for creating a user.
type CreateUser struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Role            string `json:"role"`
}

// Validate performs the validation for the CreateUser struct.
// It checks the following rules:
// - Name is required and must be between 5 and 20 characters long.
// - Email is required and must be a valid email address.
// - Phone is required and must be exactly 10 characters long.
// - Address is required and must be between 10 and 50 characters long.
// - Password is required and must be at least 10 characters long.
// - ConfirmPassword is required and must be at least 10 characters long.
// - Role is required and must be either "user" or "admin".
func (v CreateUser) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Name, validation.Required, validation.Length(5, 20)),
		validation.Field(&v.Email, validation.Required, is.Email),
		validation.Field(&v.Phone, validation.Required, validation.Length(10, 10)),
		validation.Field(&v.Address, validation.Required, validation.Length(10, 50)),
		validation.Field(&v.Password, validation.Required, validation.Length(10, 0)),
		validation.Field(&v.ConfirmPassword, validation.Required, validation.Length(10, 0)),
		validation.Field(&v.Role, validation.Required, validation.In("user", "admin")),
	)
}

// UpdateUser represents the data structure for updating a user.
type UpdateUser struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Role    string `json:"role"`
}

// Validate performs the validation for the UpdateUser struct.
// It checks the following rules:
// - Name is required and must be between 5 and 20 characters long.
// - Email is required and must be a valid email address.
// - Phone is required and must be exactly 10 characters long.
// - Address is required and must be between 10 and 50 characters long.
// - Role is required and must be either "user" or "admin".
func (v UpdateUser) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Name, validation.Required, validation.Length(5, 20)),
		validation.Field(&v.Email, validation.Required, is.Email),
		validation.Field(&v.Phone, validation.Required, validation.Length(10, 10)),
		validation.Field(&v.Address, validation.Required, validation.Length(10, 50)),
		validation.Field(&v.Role, validation.Required, validation.In("user", "admin")),
	)
}

// Login represents the data structure for user login.
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate performs the validation for the Login struct.
// It checks the following rules:
// - Email is required and must be a valid email address.
// - Password is required and must be at least 10 characters long.
func (lg Login) Validate() error {
	return validation.ValidateStruct(&lg,
		validation.Field(&lg.Email, validation.Required, is.Email),
		validation.Field(&lg.Password, validation.Required, validation.Length(10, 0)),
	)
}
