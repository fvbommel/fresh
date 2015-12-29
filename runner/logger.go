package runner

import (
	"fmt"
	logPkg "log"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
)

type logFunc func(string, ...interface{})

var logger = logPkg.New(colorable.NewColorableStderr(), "", 0)

func newLogFunc(prefix string) func(string, ...interface{}) {
	color, clear := "", ""
	if settings["colors"] == "1" {
		color = fmt.Sprintf("\033[%sm", logColor(prefix))
		clear = fmt.Sprintf("\033[%sm", colors["reset"])
	}
	prefix = fmt.Sprintf("%-11s", prefix)

	return func(format string, v ...interface{}) {
		now := time.Now()
		timeString := now.Format("15:04:03")
		format = fmt.Sprintf("%s%s %s |%s %s", color, timeString, prefix, clear, format)
		logger.Printf(format, v...)
	}
}

func fatal(err error) {
	logger.Fatal(err)
}

type appLogWriter struct{}

func (a appLogWriter) Write(p []byte) (n int, err error) {
	s := strings.TrimSuffix(string(p), "\n") // Prevent a last empty line.
	for _, line := range strings.SplitAfter(s, "\n") {
		appLog(line)
	}

	return len(p), nil
}
