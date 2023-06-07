package db

import (
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

func BuildConfig() *config {

	viper.SetConfigFile("../.env")
	viper.ReadInConfig()

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
