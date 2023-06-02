package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var DB *gorm.DB

// DBConfig represents db configuration
type Config struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func BuildConfig() *Config {

	viper.SetConfigFile("../.env")
	viper.ReadInConfig()

	dbConfig := Config{
		Host:     viper.GetString("HOST"),
		Port:     viper.GetInt("PORT"),
		User:     viper.GetString("USERNAME"),
		Password: viper.GetString("PASSWORD"),
		DBName:   viper.GetString("DBNAME"),
	}
	return &dbConfig
}
func DbURL(dbConfig *Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}
