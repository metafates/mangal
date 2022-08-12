package log

import (
	"errors"
	"fmt"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

var writeLogs bool

func Setup() error {
	writeLogs = viper.GetBool(config.LogsWrite)

	if !writeLogs {
		return nil
	}

	logsPath := where.Logs()

	if logsPath == "" {
		return errors.New("logs path is not set")
	}

	today := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(logsPath, fmt.Sprintf("%s.log", today))
	if !lo.Must(filesystem.Get().Exists(logFilePath)) {
		lo.Must(filesystem.Get().Create(logFilePath))
	}
	logFile, err := filesystem.Get().OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(logFile)

	switch viper.GetString(config.LogsLevel) {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	return nil
}

func Panic(args ...interface{}) {
	if writeLogs {
		log.Panic(args...)
	}
}

func Fatal(args ...interface{}) {
	if writeLogs {
		log.Fatal(args...)
	}
}

func Error(args ...interface{}) {
	if writeLogs {
		log.Error(args...)
	}
}

func Warn(args ...interface{}) {
	if writeLogs {
		log.Warn(args...)
	}
}

func Info(args ...interface{}) {
	if writeLogs {
		log.Info(args...)
	}
}

func Debug(args ...interface{}) {
	if writeLogs {
		log.Debug(args...)
	}
}

func Trace(args ...interface{}) {
	if writeLogs {
		log.Trace(args...)
	}
}
