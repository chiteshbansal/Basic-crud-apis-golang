package validator

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CreateUser struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Role            string `json:"role"`
}

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

type UpdateUser struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Role    string `json:"role"`
}

func (v UpdateUser) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Name, validation.Required, validation.Length(5, 20)),
		validation.Field(&v.Email, validation.Required, is.Email),
		validation.Field(&v.Phone, validation.Required, validation.Length(10, 10)),
		validation.Field(&v.Address, validation.Required, validation.Length(10, 50)),
		validation.Field(&v.Role, validation.Required, validation.In("user", "admin")),
	)
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (lg Login) Validate() error {
	return validation.ValidateStruct(&lg,
		validation.Field(&lg.Email, validation.Required, is.Email),
		validation.Field(&lg.Password, validation.Required, validation.Length(10, 0)),
	)
}
