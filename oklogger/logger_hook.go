package oklogger

import log "github.com/sirupsen/logrus"

// OkWebLogHook ：
type OkWebLogHook struct {
	Source string
}

// Fire ：
func (hook *OkWebLogHook) Fire(entry *log.Entry) error {
	if len(hook.Source) > 0 {
		entry.Data["logger"] = hook.Source
	}
	return nil
}

// Levels ：
func (hook *OkWebLogHook) Levels() []log.Level {
	return log.AllLevels
}
