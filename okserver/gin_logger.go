package okserver

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/okgotool/okgoweb/oklogger"
	"github.com/sirupsen/logrus"
)

var (
	GinLoggerConfig *oklogger.LoggerConfig = &oklogger.LoggerConfig{
		LogFileFolder:     "/var/log",
		LogFileName:       "gin.log",
		LogLevel:          logrus.InfoLevel,
		LogFileMaxSizeMb:  500,
		LogFileMaxBackups: 5,
		LogFileMaxAge:     30,
		LogFileCompress:   false,
		EnableConsole:     true,
		EnableFile:        true,
	}

	// env: GIN_MODE: debug|release|test
	GinLogLevels map[string]string = map[string]string{
		"debug":   gin.DebugMode,
		"test":    gin.TestMode,
		"info":    gin.ReleaseMode,
		"warn":    gin.ReleaseMode,
		"error":   gin.ReleaseMode,
		"release": gin.ReleaseMode,
	}
)

func InitGinLog() {

	gin.DisableConsoleColor()
	_, err := os.Create(GinLoggerConfig.LogFileFolder + "/" + GinLoggerConfig.LogFileName)
	if err == nil {
		rotateLog := oklogger.CreateRotateLogWriter(GinLoggerConfig)
		gin.DefaultWriter = io.MultiWriter(rotateLog)
	}

	logMode := gin.TestMode
	if GinLoggerConfig.LogLevel != logrus.DebugLevel {
		logMode = gin.ReleaseMode
	} else {
		logMode = gin.DebugMode
	}

	gin.SetMode(logMode)
}
