package okmysql

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dbConnectionURL string = "?charset=utf8&parseTime=True&loc=Local"
)

type (
	DbConfig struct {
		Host         string `json:"host"`
		Port         string `json:"port"`
		DbName       string `json:"dbName"`
		User         string `json:"user"`
		Password     string `json:"password"`
		MaxIdleConns int    `json:"maxIdleConns"`
		MaxOpenConns int    `json:"maxOpenConns"`
		// ConnMaxIdleTime time.Duration `json:"connMaxIdleTime"`
		// ConnMaxLifetime time.Duration `json:"connMaxLifetime"`
	}
)

// GetDBConnFromConfig :
func GetDBConnFromConfig(c *DbConfig) (*gorm.DB, error) {
	dbHost := c.Host
	dbPort := c.Port
	dbName := c.DbName
	dbUser := c.User
	dbPass := c.Password

	logger.WithFields(logrus.Fields{"dbHost": dbHost, "dbPort": dbPort, "dbName": dbName, "dbUser": dbUser}).Info("Start to init DB connection...")

	dbDSN := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + dbConnectionURL
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dbDSN, // data source name
		DefaultStringSize:         4096,  // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: NewDbLogger(),
	})

	if err != nil {
		logger.Error("Failed to connect database: "+dbUser+":@tcp("+dbHost+":"+dbPort+")/"+dbName, err)
		return nil, errors.New("failed to connect database: " + err.Error())
	}
	logger.WithFields(logrus.Fields{"dbHost": dbHost, "dbPort": dbPort, "dbName": dbName}).Info("DB connected.")

	sqlDB, err := db.DB()
	if err != nil {
		logger.WithFields(logrus.Fields{"dbHost": dbHost, "dbPort": dbPort, "dbName": dbName}).Error("Failed to get sqlDB: "+dbUser+":@tcp("+dbHost+":"+dbPort+")/"+dbName, err)
		return nil, errors.New("failed to get sqlDB: " + err.Error())
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	if c.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	} else {
		sqlDB.SetMaxIdleConns(10)
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	if c.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	} else {
		sqlDB.SetMaxOpenConns(100)
	}

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// sqlDB.SetConnMaxIdleTime(time.Hour)
	// sqlDB.SetConnMaxLifetime(24 * time.Hour)

	return db, nil

}
