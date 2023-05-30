package main

import (
	model "first-api/api/Models"
	route "first-api/api/Routes"
	"first-api/api/repository"
	"first-api/api/service"
	"first-api/api/utils"
	"fmt"

	db "first-api/database"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var err error

func main() {
	// seting up env variables
	viper.SetConfigFile("./.env")
	viper.ReadInConfig()

	userService := &service.UserService{
		Store: repository.UserStore{},
	}

	db.DB, err = gorm.Open("mysql", db.DbURL(db.BuildConfig()))
	if err != nil {
		fmt.Println("Status:", err)
	}

	defer db.DB.Close()

	db.DB.AutoMigrate(&model.User{})
	r := gin.Default()
	
	utils.RegisterRoutes(userService)
	route.InitializeRoutes(r)
	r.Run()
}
