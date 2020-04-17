package middlewares

import (
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// type Level int

// const (
// 	INFO Level = iota
// 	WARNING
// 	ERROR
// 	FATAL
// )

var (
	file *os.File
	e    error
)

// func CustomLogger() gin.HandlerFunc {
// 	gin.DisableConsoleColor()

// 	now := time.Now() //or time.Now().UTC()
// 	logFileName := "gin_" + now.Format(os.Getenv("APP_DATE_FORMAT")) + ".log"

// 	dirPath := path.Join(".", "logs")
// 	// Create directory if does not exist
// 	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
// 		err = os.MkdirAll(dirPath, os.ModeDir)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	file, e = os.OpenFile(path.Join(dirPath, logFileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
// 	if e != nil {
// 		panic(e)
// 	}

// 	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)

// 	g := gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
// 		levelFlags := []string{"INFO", "WARN", "ERROR", "FATAL"}
// 		var level string
// 		status := param.StatusCode

// 		switch {
// 		case status > 499:
// 			level = levelFlags[FATAL]
// 		case status > 399:
// 			level = levelFlags[ERROR]
// 		case status > 299:
// 			level = levelFlags[WARNING]
// 		default:
// 			level = levelFlags[INFO]
// 		}
// 		return fmt.Sprintf("[%s] - %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
// 			level,
// 			param.ClientIP,
// 			param.TimeStamp.Format(os.Getenv("APP_TIME_FORMAT")),
// 			param.Method,
// 			param.Path,
// 			param.Request.Proto,
// 			status,
// 			param.Latency,
// 			param.Request.UserAgent(),
// 			param.ErrorMessage,
// 		)
// 	})

// 	return g
// }

func CloseLogFile() {
	if err := file.Close(); err != nil {
		return
	}
}

func CustomLoggerZap() gin.HandlerFunc {

	writerSyncer := getLogWriter()
	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())

	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {

			status := c.Writer.Status()

			loggerConfig := logger.Info
			switch {
			case status > 499:
				loggerConfig = logger.Fatal
			case status > 399:
				loggerConfig = logger.Error
			case status > 299:
				loggerConfig = logger.Warn
			}

			loggerConfig(path,
				zap.Int("status", status),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("time", end.Format(os.Getenv("APP_TIME_FORMAT"))),
				zap.Duration("latency", latency),
			)
		}
	}

}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
func getLogWriter() zapcore.WriteSyncer {

	logFileName := "service.log"
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

	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(dirPath, logFileName),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
