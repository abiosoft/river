package river

import (
	oslog "log"
	"os"
)

var (
	// Log is the logger. This can be replaced or set to nil.
	Log = oslog.New(os.Stdout, "[River] ", 0)

	log riverLog
)

type riverLog struct{}

func (r riverLog) println(v ...interface{}) {
	if Log != nil {
		Log.Println(v...)
	}
}

func (r riverLog) printf(format string, v ...interface{}) {
	if Log != nil {
		Log.Printf(format, v...)
	}
}

func (r riverLog) prefix() string {
	if Log != nil {
		return Log.Prefix()
	}
	return ""
}
