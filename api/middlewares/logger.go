package middlewares

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

type Level int

const (
	INFO Level = iota
	WARNING
	ERROR
	FATAL
)

var (
	file *os.File
	e    error
)

func CustomLogger() gin.HandlerFunc {
	gin.DisableConsoleColor()

	now := time.Now() //or time.Now().UTC()
	logFileName := "gin_" + now.Format(os.Getenv("APP_DATE_FORMAT")) + ".log"

	dirPath := path.Join(".", "logs")
	// Create directory if does not exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}

	file, e = os.OpenFile(path.Join(dirPath, logFileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if e != nil {
		panic(e)
	}

	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)

	g := gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		levelFlags := []string{"INFO", "WARN", "ERROR", "FATAL"}
		var level string
		status := param.StatusCode

		switch {
		case status > 499:
			level = levelFlags[FATAL]
		case status > 399:
			level = levelFlags[ERROR]
		case status > 299:
			level = levelFlags[WARNING]
		default:
			level = levelFlags[INFO]
		}
		return fmt.Sprintf("[%s] - %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			level,
			param.ClientIP,
			param.TimeStamp.Format(os.Getenv("APP_TIME_FORMAT")),
			param.Method,
			param.Path,
			param.Request.Proto,
			status,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})

	return g
}

func CloseLogFile() {
	if err := file.Close(); err != nil {
		return
	}
}
