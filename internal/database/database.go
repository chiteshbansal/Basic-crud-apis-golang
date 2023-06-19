package db

import (
	model "first-api/internal/models"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var DB *gorm.DB

// dbConfig represents db configuration
type config struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func NewDB() {
	// Connect to the MySQL database using gorm
	var err error
	DB, err = gorm.Open("mysql", DbURL(BuildConfig()))
	if err != nil {
		fmt.Println("Status:", err)
	}

	// Automigrate user model, this will create the user table in the database
	DB.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
}

func BuildConfig() *config {

	fmt.Println(viper.GetString("DBNAME"))
	dbConfig := config{
		Host:     viper.GetString("HOST"),
		Port:     viper.GetInt("PORT"),
		User:     viper.GetString("USERNAME"),
		Password: viper.GetString("PASSWORD"),
		DBName:   viper.GetString("DBNAME"),
	}
	return &dbConfig
}

func DbURL(dbConfig *config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}
