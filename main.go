package main

import(
	"first-api/Config"
	"first-api/Routes"
	"first-api/Models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

var err error

func main(){
	Config.DB , err = gorm.Open("mysql",Config.DbURL(Config.BuildDBConfig()))
	if err !=nil{
		fmt.Println("Status:",err)
	}

	defer Config.DB.Close()
	Config.DB.AutoMigrate(&Models.User{})
	r := Routes.SetupRouter()
	r.Run()
}