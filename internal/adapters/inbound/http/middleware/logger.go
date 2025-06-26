package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)

		stop := time.Now()

		method := c.Request().Method
		path := c.Request().URL.Path
		status := c.Response().Status
		latency := stop.Sub(start)

		log.Printf("[%d] %s %s (%s)", status, method, path, latency)
		return err
	}
}
