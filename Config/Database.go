package config

import (
	"fmt"

	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var DB *gorm.DB

// DBConfig represents db configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func BuildDBConfig() *DBConfig {
	viper.SetConfigFile("../.env")
	viper.ReadInConfig()
	port, _ := strconv.Atoi(viper.Get("PORT").(string))
	dbConfig := DBConfig{
		Host:     "localhost",
		Port:     port,
		User:     "root",
		Password: "pass123",
		DBName:   "firstapi",
	}

	return &dbConfig
}
func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}
