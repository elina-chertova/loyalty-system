// Package logger provides logging functionalities for the application.
// It utilizes the zap logging library to log HTTP requests and other
// important information throughout the application.
package logger

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a global zap.Logger instance used throughout the application.
// It is initialized to a no-operation logger by default and should be
// initialized with InitLogger at the start of the application.
var Logger *zap.Logger = zap.NewNop()

// InitLogger initializes the global Logger with a production configuration.
// It sets up structured, leveled logging using the zap library.
func InitLogger() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger
	return nil
}

// GinLogger returns a gin.HandlerFunc (middleware) that logs requests.
// It logs various details about each HTTP request including method, endpoint,
// status code, duration, client IP, and user agent. It handles different
// logging details based on the HTTP method (POST or GET).
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		var size any
		if c.Request.Method == http.MethodPost {
			size = c.Request.ContentLength
		} else if c.Request.Method == http.MethodGet {
			size = c.Writer.Size()
		}

		logger.Info(
			"Request",
			zap.String("method", c.Request.Method),
			zap.String("endpoint", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Any("body_size", size),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}
