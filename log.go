package river

import (
	"fmt"
	oslog "log"
	"os"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

var (
	// Log is the logger. This can be replaced or set to nil.
	Log = oslog.New(os.Stdout, "[River] ", 0)

	// LogRequests enables request log. Useful for development.
	LogRequests = true

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

// requestLogger is a middleware that logs requests in a colourful way.
func requestLogger() Middleware {
	return func(c *Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		bg := color.BgBlack
		switch {
		case c.Status() >= 200 && c.Status() < 300:
			bg = color.BgGreen
		case c.Status() >= 300 && c.Status() < 400:
			bg = color.BgBlue
		case c.Status() >= 400 && c.Status() < 500:
			bg = color.BgYellow
		case c.Status() >= 500 && c.Status() < 600:
			bg = color.BgRed
		}

		paint := color.New(bg, color.FgWhite, color.Bold).SprintFunc()
		status := paint(fmt.Sprintf("  %d  ", c.Status()))
		size := humanize.Bytes(uint64(c.Written()))

		fmt.Printf("%s%v|%s|%15v|%6s|%-4s %s\n",
			log.prefix(),
			time.Now().Format("2006-01-02 15:04:05"),
			status, duration, size, c.Method, c.URL.Path,
		)

	}
}
