package okmysql

import (
	"errors"

	"gorm.io/gorm"
)

var (
	dbConnConfigs map[string]*DbConfig = map[string]*DbConfig{}
	dbConnPool    map[string]*gorm.DB  = map[string]*gorm.DB{}
)

func AddConnConfig(name string, config *DbConfig) error {
	if cfg, ok := dbConnConfigs[name]; ok && cfg == config {
		if conn, ok := dbConnPool[name]; ok && conn != nil {
			return nil
		}
	}

	conn, err := GetDBConnFromConfig(config)
	if err != nil || conn == nil {
		return errors.New("Failed to connect db with the config, invalid db config:" + err.Error())
	}

	dbConnConfigs[name] = config

	// reset connection:
	dbConnPool[name] = conn

	return nil
}

func GetDbConn(connName string) (*gorm.DB, error) {
	if conn, ok := dbConnPool[connName]; ok && conn != nil {
		return conn, nil
	} else if config, ok := dbConnConfigs[connName]; ok && config != nil {
		conn, err := GetDBConnFromConfig(config)
		if err == nil && conn != nil {
			dbConnPool[connName] = conn
		}
		return conn, err
	} else {
		return nil, errors.New("Not found db connection config for name:" + connName)
	}
}
