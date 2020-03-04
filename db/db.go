package db

import (
	"flow/logger"
	"fmt"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db *gorm.DB
var err error

type DBService struct{}

// Init : Initializes the database migrations
func Init() {
	dbUserName := viper.GetString("database.username")
	dbPassword := viper.GetString("database.password")
	dbURL := viper.GetString("database.url")
	dbName := viper.GetString("database.name")
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbURL, dbUserName, dbName, dbPassword) //Build connection string

	//dbConnectionString := dbUserName + ":" + dbPassword + "@tcp(" + dbUrl + ")/" + dbName
	db, err = gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Println("failed to connect.", dbURI, err)
		logger.SugarLogger.Info("Failed to connect to DB", dbURI, err.Error())
	}
}

// GetDB : Get an instance of DB to connect to the database connection pool
func (d DBService) GetDB() *gorm.DB {
	return db
}
