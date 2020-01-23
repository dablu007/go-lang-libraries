package db

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"fmt"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/lib/pq"
	logrs "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
)

var db *gorm.DB
var err error

func Init() {
	dbUserName := viper.GetString("database.username")
	dbPassword := viper.GetString("database.password")
	dbUrl := viper.GetString("database.url")
	dbName := viper.GetString("database.name")
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbUrl, dbUserName, dbName, dbPassword) //Build connection string

	//dbConnectionString := dbUserName + ":" + dbPassword + "@tcp(" + dbUrl + ")/" + dbName
	db, err = gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Println("failed to connect.", dbUri, err)
	}
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(" Not able to fetch the working directory")
		return
	}
	db.SingularTable(true)
	workingDir = workingDir + "/db/migrations"
	migrateConf := &goose.DBConf{
		MigrationsDir: workingDir,
		Driver: goose.DBDriver{
			Name:    "postgres",
			OpenStr: dbUri,
			Import:  "github.com/lib/pq",
			Dialect: &goose.PostgresDialect{},
		},
	}
	logrs.Println(" Fetching the most recent DB version ")
	latest, err := goose.GetMostRecentDBVersion(migrateConf.MigrationsDir)
	if err != nil {
		log.Println(err)

	}
	fmt.Println(" Most recent DB version ", latest)
	logrs.Println("Running the migrations on db", workingDir)
	err = goose.RunMigrationsOnDb(migrateConf, migrateConf.MigrationsDir, latest, db.DB())
	if err != nil {
		log.Println(err)
	}
}

func GetDB() *gorm.DB {
	return db
}
