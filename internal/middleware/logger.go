package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		// Safely get request ID, default to empty string if not set
        requestID := ""
        if id := c.Locals("requestid"); id != nil {
            if strID, ok := id.(string); ok {
                requestID = strID
            }
        }

		logger.Info("Request",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("ip", c.IP()),
			zap.Duration("latency", time.Since(start)),
			zap.String("user_agent", c.Get("User-Agent")),
		)

		return err
	}
}
