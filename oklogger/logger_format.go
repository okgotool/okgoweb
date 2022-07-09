package oklogger

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type LogFormatter struct{}

//格式详情
func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timeStr := time.Now().Local().Format("2006-01-02 15:04:05.000")
	timestamp := strings.Replace(timeStr, ".", ",", 1)
	var file string
	var leng int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		leng = entry.Caller.Line
	}
	var msg string
	if entry.Message != "" {
		msg = entry.Message
	}

	if len(entry.Data) > 0 {
		field := toString(entry.Data)
		msg = fmt.Sprintf(msg+" %s", field)
	}
	var level string
	if entry.Level != log.WarnLevel {
		level = strings.ToUpper(entry.Level.String())
	} else {
		level = "WARN"
	}
	message := fmt.Sprintf("%s %s [%s:%d] - level=%s msg=%s\n", timestamp, level, file, leng, level, msg)
	return []byte(message), nil
}

func toString(fields log.Fields) string {
	tmp := make([]string, 0, len(fields))

	for k, v := range fields {
		field := fmt.Sprintf(" %s=%s", k, v)
		tmp = append(tmp, field)
	}

	return strings.Join(tmp, " ")
}
