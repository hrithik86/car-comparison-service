package logger

import (
	"car-comparison-service/config"
	"context"
	joonix "github.com/joonix/log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	AppType = "appType"
)

type Logger struct {
	*logrus.Logger
}

var Log *Logger

type LogError struct {
	Error error
}

func panicIfError(err error) {
	if err != nil {
		panic(LogError{err})
	}
}

func SetupLogger() {
	level, err := logrus.ParseLevel(config.LogLevel())
	panicIfError(err)
	var formatter = joonix.NewFormatter()
	formatter.TimestampFormat = func(fields logrus.Fields, now time.Time) error {
		fields["timeStamp"] = now.Format(time.RFC3339)
		return nil
	}

	logrusVar := &logrus.Logger{
		Out:          os.Stdout,
		Formatter:    formatter,
		Hooks:        make(logrus.LevelHooks),
		Level:        level,
		ReportCaller: true,
	}

	Log = &Logger{logrusVar}
}

func Get(ctx context.Context) *logrus.Entry {
	entry := Log.WithField("serviceName", "car-comparison-service")
	if appType := ctx.Value(AppType); appType != nil {
		entry = entry.WithField(AppType, appType)
	}
	return entry
}
