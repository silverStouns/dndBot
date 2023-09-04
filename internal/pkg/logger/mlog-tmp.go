package logger

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"path"
)

const (
	prefixPath = ""
)

var (
	logger *logs.BeeLogger
)

//goland:noinspection GoLinter
func init() {
	var err error
	filename := "BOT" //filepath.Base(os.Args[0])
	fullPath := "./" + path.Join(prefixPath, filename, fmt.Sprint(filename, "-time.log"))

	lc := `{"filename":"` + fullPath + `", "level":7, "maxlines":0, "maxsize":5000000, "daily":true, "maxdays":14, "color":true, "rotate":true }`

	logger = logs.NewLogger(5000)
	// Включаем отображение файла и номер строки
	logger.EnableFuncCallDepth(true)
	// Поднимаем выше уровень выше для EnableFuncCallDepth, чтобы не отображалось [logger-tmp:30]
	logger.SetLogFuncCallDepth(3)

	if err = logger.SetLogger(logs.AdapterConsole, lc); err != nil {
		panic(err)
	}

	if err = logger.SetLogger(logs.AdapterFile, lc); err != nil {
		panic(err)
	}
}

// TraceManyMessage ...
//
//goland:noinspection ALL
func TraceManyMessage(format string, v ...interface{}) {
	spf := fmt.Sprintf(format, v...)
	if len(spf) <= 512 {
		Trace(spf)
	} else {
		Trace(spf[:512])
	}
}

// Trace ...
func Trace(format string, v ...interface{}) {
	logger.Trace(format, v...)
}

// Debug ...
//
//goland:noinspection ALL
func Debug(format string, v ...interface{}) {
	logger.Debug(format, v...)
}

// Notice ...
//
//goland:noinspection ALL
func Notice(format string, v ...interface{}) {
	logger.Notice(format, v...)
}

// Info ...
//
//goland:noinspection ALL
func Info(format string, v ...interface{}) {
	logger.Info(format, v...)
}

// Alert ...
//
//goland:noinspection ALL
func Alert(format string, v ...interface{}) {
	logger.Alert(format, v...)
}

// Warn ...
//
//goland:noinspection ALL
func Warn(format string, v ...interface{}) {
	logger.Warn(format, v...)
}

// Error ...
func Error(format string, v ...interface{}) {
	logger.Error(format, v...)
}

// Critical ...
func Critical(format string, v ...interface{}) {
	logger.Critical(format, v...)
}
