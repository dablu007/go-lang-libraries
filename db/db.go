package db

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"fmt"
	"github.com/dablu007/go-lang-libraries/logger"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"os"
	"time"
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
	maxIdleConnections := viper.GetInt("postgresdb.maxIdleConnections")
	maxOpenConnections := viper.GetInt("postgresdb.maxOpenConns")
	connectionMaxLifetime := viper.GetInt("postgresdb.connMaxLifetimeInHours")

	//dbConnectionString := dbUserName + ":" + dbPassword + "@tcp(" + dbUrl + ")/" + dbName
	db, err = gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Println("failed to connect.", dbURI, err)
		logger.SugarLogger.Fatalf("Failed to connect to DB", dbURI, err.Error())
		os.Exit(1)
	}
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Not able to fetch the working directory")
		logger.SugarLogger.Fatalf("Not able to fetch the working directory")
		os.Exit(1)
	}
	db.DB().SetMaxIdleConns(maxIdleConnections)
	db.DB().SetMaxOpenConns(maxOpenConnections)
	db.DB().SetConnMaxLifetime(time.Hour * time.Duration(connectionMaxLifetime))
	db.SingularTable(true)
	workingDir = workingDir + "/db/migrations"
	migrateConf := &goose.DBConf{
		MigrationsDir: workingDir,
		Driver: goose.DBDriver{
			Name:    "postgres",
			OpenStr: dbURI,
			Import:  "github.com/lib/pq",
			Dialect: &goose.PostgresDialect{},
		},
	}
	logger.SugarLogger.Infof("Fetching the most recent DB version")
	latest, err := goose.GetMostRecentDBVersion(migrateConf.MigrationsDir)
	if err != nil {
		logger.SugarLogger.Errorf("Unable to get recent goose db version", err)

	}
	fmt.Println(" Most recent DB version ", latest)
	logger.SugarLogger.Infof("Running the migrations on db", workingDir)
	err = goose.RunMigrationsOnDb(migrateConf, migrateConf.MigrationsDir, latest, db.DB())
	if err != nil {
		logger.SugarLogger.Fatalf("Error while running migrations", err)
		os.Exit(1)
	}
}

// GetDB : Get an instance of DB to connect to the database connection pool
func (d DBService) GetDB() *gorm.DB {
	return db
}
