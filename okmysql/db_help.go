package okmysql

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// const (
// 	dbConnectionURL string = "?charset=utf8&parseTime=True&loc=Local"
// )

type (
	DbConfig struct {
		Host         string `json:"host" yaml:"host"`
		Port         int    `json:"port" yaml:"port"`
		DbName       string `json:"dbName" yaml:"dbName"`
		User         string `json:"user" yaml:"user"`
		Password     string `json:"password" yaml:"password"`
		MaxIdleConns int    `json:"maxIdleConns" yaml:"maxIdleConns"`
		MaxOpenConns int    `json:"maxOpenConns" yaml:"maxOpenConns"`
		Dsn          string `json:"dsn" yaml:"dsn"`
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

	if len(dbHost) < 1 || len(dbName) < 1 || len(dbPass) < 1 {
		return nil, errors.New("DB host or name or password not set")
	}

	if dbPort < 1 {
		dbPort = 3306
	}
	if len(dbUser) < 1 {
		dbUser = "root"
	}

	logger.WithFields(logrus.Fields{"dbHost": dbHost, "dbPort": dbPort, "dbName": dbName, "dbUser": dbUser}).Info("Start to init DB connection...")

	dbDSN := fmt.Sprintf(dbUser+":"+dbPass+"@tcp("+dbHost+":%d)/"+dbName+"?"+c.Dsn, dbPort)
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
		logger.Error("Failed to connect database: "+dbUser+":@tcp("+dbHost+":", dbPort, ")/"+dbName, err)
		return nil, errors.New("failed to connect database: " + err.Error())
	}
	logger.WithFields(logrus.Fields{"dbHost": dbHost, "dbPort": dbPort, "dbName": dbName}).Info("DB connected.")

	sqlDB, err := db.DB()
	if err != nil {
		logger.WithFields(logrus.Fields{"dbHost": dbHost, "dbPort": dbPort, "dbName": dbName}).Error("Failed to get sqlDB: "+dbUser+":@tcp("+dbHost+":", dbPort, ")/"+dbName, err)
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
