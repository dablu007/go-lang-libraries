package db

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"flow/config"
	"fmt"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/lib/pq"
	logrs "github.com/sirupsen/logrus"
	"log"
	"os"
)

var db *gorm.DB
var err error

func Init() {
	config := config.GetConfig()
	dbUserName := config.GetString("database.username")
	dbPassword := config.GetString("database.password")
	dbUrl := config.GetString("database.url")
	dbName := config.GetString("database.name")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbUrl, dbUserName, dbName, dbPassword) //Build connection string

	//dbConnectionString := dbUserName + ":" + dbPassword + "@tcp(" + dbUrl + ")/" + dbName
	db, err = gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Println("failed to connect.", dbUri)
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
