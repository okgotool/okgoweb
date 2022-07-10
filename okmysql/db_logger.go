package okmysql

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/okgotool/okgoweb/oklogger"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var (
	DbLoggerConfig *oklogger.LoggerConfig = &oklogger.LoggerConfig{
		LogFileFolder:     "/var/log",
		LogFileName:       "db.log",
		LogLevel:          logrus.InfoLevel,
		LogFileMaxSizeMb:  500,
		LogFileMaxBackups: 5,
		LogFileMaxAge:     30,
		LogFileCompress:   false,
		EnableConsole:     true,
		EnableFile:        true,
	}
)

type dblogger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	logWiter              *logrus.Logger
}

func NewDbLogger() *dblogger {
	logger := oklogger.CreateNewLogger("", DbLoggerConfig)

	return &dblogger{
		logWiter:              logger,
		SkipErrRecordNotFound: true,
		SlowThreshold:         500 * time.Millisecond, // 慢sql为500ms
	}
}

func (l *dblogger) SetLevel(level string) {
	logrusLevel, ok := oklogger.AppLogLevels[strings.ToLower(level)]
	if ok {
		l.logWiter.SetLevel(logrusLevel)
	}
}

func (l *dblogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *dblogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.logWiter.WithContext(ctx).Infof(s, args)
}

func (l *dblogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.logWiter.WithContext(ctx).Warnf(s, args)
}

func (l *dblogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.logWiter.WithContext(ctx).Errorf(s, args)
}

func (l *dblogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := log.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[log.ErrorKey] = err
		l.logWiter.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logWiter.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	l.logWiter.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
}
