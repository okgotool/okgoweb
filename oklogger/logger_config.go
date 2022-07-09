package oklogger

import "github.com/sirupsen/logrus"

var (
	DefaultLoggerConfig *LoggerConfig = &LoggerConfig{
		LogFileFolder:     "/var/log",
		LogFileName:       "app.log",
		LogLevel:          logrus.InfoLevel,
		LogFileMaxSizeMb:  500,
		LogFileMaxBackups: 5,
		LogFileMaxAge:     30,
		LogFileCompress:   false,
		EnableConsole:     true,
		EnableFile:        true,
	}
)

type (
	LoggerConfig struct {
		LogFileFolder     string
		LogFileName       string
		LogLevel          logrus.Level
		LogFileMaxSizeMb  int
		LogFileMaxBackups int
		LogFileMaxAge     int
		LogFileCompress   bool
		EnableConsole     bool
		EnableFile        bool
	}
)

func NewLoggerConfig(c *LoggerConfig) *LoggerConfig {
	if len(c.LogFileFolder) < 1 {
		c.LogFileFolder = DefaultLoggerConfig.LogFileFolder
	}
	if len(c.LogFileName) < 1 {
		c.LogFileName = DefaultLoggerConfig.LogFileName
	}
	if c.LogLevel == 0 {
		c.LogLevel = DefaultLoggerConfig.LogLevel
	}
	if c.LogFileMaxSizeMb == 0 {
		c.LogFileMaxSizeMb = DefaultLoggerConfig.LogFileMaxSizeMb
	}
	if c.LogFileMaxBackups == 0 {
		c.LogFileMaxBackups = DefaultLoggerConfig.LogFileMaxBackups
	}
	if c.LogFileMaxAge == 0 {
		c.LogFileMaxAge = DefaultLoggerConfig.LogFileMaxAge
	}

	return c
}
