package main

import (
	config "first-api/config"
	model "first-api/Models"
	route "first-api/Routes"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var err error

func main() {
	// seting up env variables
	viper.SetConfigFile("./.env")
	viper.ReadInConfig()

	config.DB, err = gorm.Open("mysql", config.DbURL(config.BuildDBConfig()))
	if err != nil {
		fmt.Println("Status:", err)
	}

	userStore := new(model.UserStore)
	defer config.DB.Close()
	config.DB.AutoMigrate(&model.User{})
	r := route.SetupRouter(userStore)
	r.Run()
}
