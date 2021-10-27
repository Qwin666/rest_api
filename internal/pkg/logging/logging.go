package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var e *logrus.Entry

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (h *writerHook) Fire(e *logrus.Entry) error {
	line, err := e.String()
	if err != nil {
		return err
	}
	for _, w := range h.Writer {
		w.Write([]byte(line))
	}
	return err
}

func (h *writerHook) Levels() []logrus.Level {
	return h.LogLevels
}

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	return &Logger{e}
}

func (l *Logger) GetLogerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

func init() {
	i := logrus.New()
	i.SetReportCaller(true)
	i.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fileName := path.Base(f.File)
			return fmt.Sprintf("%s", f.Function), fmt.Sprintf("%s:%d\n", fileName, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}
	err := os.MkdirAll("logs", 0644)
	if err != nil {
		panic(err)
	}
	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}
	
	i.SetOutput(io.Discard)

	i.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})
	i.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(i)
}
