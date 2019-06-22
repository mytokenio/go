package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type LineHook struct {
	debug     bool
	Field     string
	Skip      int
	Formatter func(file, function string, line int) string
}

func NewLineHook(debug bool) *LineHook {
	return &LineHook{
		debug: debug,
		Field: "[Fn]",
		Skip:  5,
		Formatter: func(file, function string, line int) string {
			return fmt.Sprintf("%s:%d", file, line)
		},
	}
}

func (hook *LineHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = hook.Formatter(findCaller(hook.Skip))
	return nil
}

func (hook *LineHook) Levels() []logrus.Level {
	if hook.debug {
		return logrus.AllLevels
	} else {
		return []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func findCaller(skip int) (string, string, int) {
	var (
		pc       uintptr
		file     string
		function string
		line     int
	)
	for i := 0; i < 10; i++ {
		pc, file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") && !strings.HasPrefix(file, "log") {
			break
		}
	}
	if pc != 0 {
		frames := runtime.CallersFrames([]uintptr{pc})
		frame, _ := frames.Next()
		function = frame.Function
	}

	return file, function, line
}

func getCaller(skip int) (uintptr, string, int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return 0, "", 0
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n += 1
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}

	return pc, file, line
}
