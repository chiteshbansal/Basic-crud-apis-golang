package Models

import(
	"fmt"
	"first-api/Config"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)




type UserStore struct{}

func (u *UserStore) CreateUser(user *User) error{	
	if err := Config.DB.Create(user).Error ; err!=nil{
		return err
	}
	return nil
}

func (us *UserStore) Validate(u User) error {
	return  validation.ValidateStruct(&u,validation.Field(&u.Name,validation.Required,validation.Length(5,20)),
	validation.Field(&u.Email,validation.Required,is.Email)	,
	validation.Field(&u.Phone, validation.Required,validation.Length(10,10)),
	validation.Field(&u.Address,validation.Required,validation.Length(10,50)),
	)

}


// get all users Fetch all user data
 func GetAllUsers(user *[]User) (err error){
	if err = Config.DB.Find(user).Error; err !=nil{
		return err
	}
	return nil
 }

 // CreateUser

 func CreateUser(user *User) (err error) {
	if err = Config.DB.Create(user).Error ; err!=nil{
		return err
	}
	return nil
 }

 // getuserById

 func GetUserByID(user *User, id string) (err error){
	if err = Config.DB.Where("id = ?",id).First(user).Error; err !=nil{
		return err
	}
	return nil
 }

 // update user 

 func UpdateUser(user *User,id string) (err error){
	fmt.Println(user)
 Config.DB.Save(user)
	return nil
 }

 // Delete User

 func DeleteUser(user *User,id string) (err error){
	fmt.Println("user: ", user)
	if err = Config.DB.Where("id = ?",id).Delete(user).Error ; err!=nil{
		return err
	}
	return nil

 }