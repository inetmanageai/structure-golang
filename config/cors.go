package config

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins: Env.Cors,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}
}
