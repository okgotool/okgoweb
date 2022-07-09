package oklogger

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (

	// Logger :
	Logger *logrus.Logger = CreateDefaultAppLogger()
	logger *logrus.Logger = Logger

	// ValidLogLevelMap :
	AppLogLevels map[string]logrus.Level = map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
	}

	// // env: GIN_MODE: debug|release|test
	// GinLogLevels map[string]string = map[string]string{
	// 	"debug":   gin.DebugMode,
	// 	"test":    gin.TestMode,
	// 	"info":    gin.ReleaseMode,
	// 	"warn":    gin.ReleaseMode,
	// 	"error":   gin.ReleaseMode,
	// 	"release": gin.ReleaseMode,
	// }

)

// // InitLogger :
// func InitLogger() {
// 	// Logger.SetLevel(defaultAppLogLevel)

// 	// reset app logger:
// 	initAppLogger()

// 	// reset gin logger:
// 	// initGinLog()
// }

// func initGinLog() {

// 	// gin.DisableConsoleColor()
// 	// file, err := os.Create(logFilePath + "/" + ginLogFileName)
// 	// if err == nil {

// 	// if !strings.EqualFold(LoggerEnv, ENV_LOCAL) {

// 	// 	logFile := logFilePath + "/" + GinLogFileName
// 	// 	prepareLogFile(logFile)

// 	// 	rotateLog := createRotateLogWriter(logFile, GinLogFileMaxBackups)
// 	// 	gin.DefaultWriter = io.MultiWriter(rotateLog)
// 	// }
// 	// }

// 	// logMode := gin.TestMode
// 	// if strings.EqualFold(loggerEnv, ENV_PROD) {
// 	// 	logMode = gin.ReleaseMode
// 	// } else if strings.EqualFold(loggerEnv, ENV_LOCAL) {
// 	// 	logMode = gin.DebugMode
// 	// }

// 	// Logger.WithFields(log.Fields{"logLevel": logMode}).Info("Init Gin log")
// 	// gin.SetMode(defaultGinLogLevel)
// }

// GetLogger :
func CreateDefaultAppLogger() *logrus.Logger {
	return CreateNewLogger("", DefaultLoggerConfig)
}

// CreateNewLogger :
func CreateNewLogger(name string, cfg *LoggerConfig) *logrus.Logger {

	lg := logrus.New()
	lg.AddHook(&OkWebLogHook{Source: name})

	config := NewLoggerConfig(cfg)
	ResetLogger(lg, config)

	return lg
}

func ResetLogger(lg *log.Logger, cfg *LoggerConfig) {

	// create log file if not local env:
	resetLoggerOutput(lg, cfg)

	// log format: TextFormatter | JSONFormatter
	lg.SetFormatter(new(LogFormatter))

	// add caller as 'method' in log: true(will add measurable overhead) | false
	lg.SetReportCaller(true)

	// default(prod stage env):
	// logrusLevel := GetAppLogLevel()

	lg.SetLevel(cfg.LogLevel)
}

func resetLoggerOutput(lg *log.Logger, cfg *LoggerConfig) {
	cfg = NewLoggerConfig(cfg)

	writers := []io.Writer{}
	if cfg.EnableConsole {
		writers = append(writers, os.Stderr)
	}
	if cfg.EnableFile {
		// logFile := cfg.LogFileFolder + "/" + cfg.LogFileName
		prepareLogFolder(cfg.LogFileFolder)

		rotateLog := CreateRotateLogWriter(cfg)
		writers = append(writers, rotateLog)
	}

	if len(writers) < 1 {
		writers = append(writers, os.Stderr)
	}
	wrt := io.MultiWriter(writers...)

	lg.SetOutput(wrt)
}

func CreateRotateLogWriter(cfg *LoggerConfig) *lumberjack.Logger {
	prepareLogFolder(cfg.LogFileFolder)

	rotateLog := &lumberjack.Logger{
		// 日志输出文件路径
		Filename: cfg.LogFileFolder + "/" + cfg.LogFileName,
		// 日志文件最大 size, 单位是 MB
		MaxSize: cfg.LogFileMaxSizeMb, // megabytes
		// 最大过期日志保留的个数
		MaxBackups: cfg.LogFileMaxBackups,
		// 保留过期文件的最大时间间隔,单位是天
		MaxAge: cfg.LogFileMaxAge, //days
		// 是否需要压缩滚动日志, 使用的 gzip 压缩
		Compress: cfg.LogFileCompress, // disabled by default
	}

	return rotateLog
}

func prepareLogFolder(fileFolder string) {

	err := os.MkdirAll(fileFolder, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func ResetLogLevel(lg *logrus.Logger, level string) {
	logrusLevel, ok := AppLogLevels[strings.ToLower(level)]
	if ok {
		lg.SetLevel(logrusLevel)
	}
}

// // GetSource :
// func GetSource(up uintptr, appName string, file string, line int, ok bool) string {
// 	source := "system"
// 	split := "/" + appName + "/"

// 	strs := strings.Split(file, split)
// 	if len(strs) > 1 {
// 		source = strs[1]
// 	}

// 	return fmt.Sprintf("%s:%d", source, line)
// }
