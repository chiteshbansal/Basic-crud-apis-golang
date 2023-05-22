package Models

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)
type User struct{

	Id uint `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Address string `json:"address"`
}

func (b *User) TableName() string{
	return "user"
}

type UserData struct {
	
	Id uint `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Address string `json:"address"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,validation.Field(&u.Name,validation.Required,validation.Length(5,20)),
	validation.Field(&u.Email,validation.Required,is.Email)	,
	validation.Field(&u.Phone, validation.Required,validation.Length(10,10)),
	validation.Field(&u.Address,validation.Required,validation.Length(10,50)),
	)

}