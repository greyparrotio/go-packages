package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *lumberjack.Logger

	// Sets the number of days for which log files will be kept.
	// The default is that 1 day of log files will be kept.
	LogRetentionDuration = 1
)

func init() {
	Setup(LogRetentionDuration)
	loggerCron := cron.New()
	loggerCron.AddFunc("@midnight", func() {
		log.Println("resetting logging destination ...")
		Setup(LogRetentionDuration)
	})
	loggerCron.Start()
}

func Setup(logRetentionDuration int) {
	LogRetentionDuration = logRetentionDuration
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	var logger io.Writer
	logsDirectory := "logs/"
	fileName := logsDirectory + "app-" + time.Now().Format(time.DateOnly) + ".log"
	_, err := os.Stat(logsDirectory)
	if err != nil {
		logger = os.Stdout
	} else {
		logFiles, err := os.ReadDir(logsDirectory)
		if err != nil {
			log.Println(err)
		}
		var validLogFileNames []string
		for duration := range LogRetentionDuration {
			validLogFileNames = append(validLogFileNames,
				logsDirectory+"app-"+
					time.Now().Add(-time.Duration(duration*24)*time.Hour).
						Format(time.DateOnly)+
					".log")
		}
		for _, logFile := range logFiles {
			if !slices.Contains(validLogFileNames, logsDirectory+logFile.Name()) {
				err = os.Remove(logsDirectory + logFile.Name())
				if err != nil {
					log.Println(err)
				}
			}
		}
		logger = &lumberjack.Logger{
			Filename:  fileName,
			Compress:  false,
			LocalTime: true,
		}
	}
	log.SetOutput(logger)
}

/*
Logger is a http handler decorator which logs information that is handled by the decorated
http handler. Logger logs this information to the command line.

The log output is in this order:

	The http method used
	The URI of the request
	The name of the http handler that handled the request
	The time spent by the handler to service the request
	The ip address(s) of the source
*/
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)

		elapsedTime := time.Since(start)

		message := fmt.Sprintf(
			"%s\t%s\t%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			elapsedTime,
			r.Header.Get("X-Forwarded-For"),
			r.RemoteAddr,
		)
		log.Println(message)
	})
}
