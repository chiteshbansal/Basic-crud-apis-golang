package controller

// package name should start with small letters
import (
	"first-api/Models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// get all users

type UserStorer interface {
	CreateUser(user *model.User) error
	Validate(user model.User) error
	GetAllUsers(users *[]model.User) error
	GetUserByID(user *model.User, id string) error
	UpdateUser(user *model.User, id string) error
	DeleteUser(user *model.User, id string) error
}

func GetUserByIDController(store UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		var user model.User

		err := store.GetUserByID(&user, id)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, user)
		}
	}
}

func GetUsers(store UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []model.User
		err := store.GetAllUsers(&users)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, users)
		}

	}
}

func NewUserController(store UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		c.BindJSON(&user)

		valErr := store.Validate(user)
		if valErr != nil {
			fmt.Println(valErr)
			c.JSON(http.StatusNotFound, valErr)
			return
		}

		err := store.CreateUser(&user)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			fmt.Println(user.Name)
			c.JSON(http.StatusOK, user)
		}
	}
}

// update user data
func UpdateUser(store UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		var user model.User
		err := store.GetUserByID(&user, id)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.BindJSON(&user) // putting the user data from the body to the user variable
		valErr := user.Validate()
		if valErr != nil {
			fmt.Println(valErr)
			c.JSON(http.StatusNotFound, valErr)
			return
		}
		err = store.UpdateUser(&user, id)

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, user)
		}

	}
}

// delete user

func DeleteUser(store UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		id := c.Params.ByName("id")

		err := store.GetUserByID(&user, id)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusNotFound, "REcord not found")
			return
		}

		c.BindJSON(&user)
		err = store.DeleteUser(&user, id)

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, gin.H{"id" + id: "is deleted"})
		}
	}

}
