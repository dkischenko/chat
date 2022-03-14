package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

type Logger struct {
	Entry *logrus.Entry
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}

	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func GetLogger() (logger *Logger, err error) {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: true,
		FullTimestamp: true,
	}

	err = os.MkdirAll(fmt.Sprintf("%s/logs", os.Getenv("CHATLOGPATH")), 0755)
	if err != nil {
		return nil, err
	}
	allFile, err := os.OpenFile(fmt.Sprintf("%s/logs/all.log", os.Getenv("CHATLOGPATH")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)

	if err != nil {
		return nil, err
	}

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stderr},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	return &Logger{Entry: logrus.NewEntry(l)}, nil
}
