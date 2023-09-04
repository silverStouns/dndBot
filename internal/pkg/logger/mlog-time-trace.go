package logger

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"path"
)

var (
	loggerTime *logs.BeeLogger
)

//goland:noinspection ALL
func init() {
	var err error
	filename := "BOT" //filepath.Base(os.Args[0])
	fullPath := "./" + path.Join(prefixPath, filename, fmt.Sprint(filename, "-time.log"))

	lc := `{"filename":"` + fullPath + `", "level":7, "maxlines":0, "maxsize":5000000, "daily":true, "maxdays":14, "color":true, "rotate":true }`

	loggerTime = logs.NewLogger(5000)
	// Включаем отображение файла и номер строки
	loggerTime.EnableFuncCallDepth(true)
	// Поднимаем выше уровень выше для EnableFuncCallDepth, чтобы не отображалось [logger-tmp:30]
	loggerTime.SetLogFuncCallDepth(3)

	if err = loggerTime.SetLogger(logs.AdapterConsole, lc); err != nil {
		panic(err)
	}

	if err = loggerTime.SetLogger(logs.AdapterFile, lc); err != nil {
		panic(err)
	}
}

// TraceTime ...
//
//goland:noinspection ALL
func TraceTime(format string, v ...interface{}) {
	loggerTime.Trace(format, v...)
}
