package main

import (
	"structure-golang/common/events"
	"structure-golang/common/logs"
	"structure-golang/config"
	"structure-golang/core/handlers"
	"structure-golang/core/middlewares"
	"structure-golang/core/repositories"
	"structure-golang/core/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func init() {
	config.NewAppInitEnvironment()
}

func main() {
	// Init commons
	event := events.NewEventKafka()
	db := config.NewAppDatabase()
	log := logs.NewAppLogsElk(config.NewAppElastic())

	// Create a new Fiber instance
	app := fiber.New()

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(config.CorsConfig()))

	// Repositories
	userRepo := repositories.NewUserRepository(db, "users")

	// Services
	userSrv := services.NewUserService(log, userRepo)

	// Handlers
	userHand := handlers.NewUserHandler(userSrv)
	consumerHand := handlers.NewConsumerHandler(log, userSrv)

	// Start Routing ---------------

	// User
	app.Post("/api/v1/signin", userHand.Signin)

	// Not Allow Method
	app.Use("*", middlewares.UnknowMethod)

	// Listen servers
	env := strings.ToLower(config.Env.Env)

	if env == "dev" || env == "development" {
		go app.Listen("localhost:" + config.Env.Port)
	} else {
		go app.Listen(":" + config.Env.Port)
	}

	// Start Event Server with Kafka Consumer ---------------
	go event.On("example_topic", config.Env.KafkaConsumerGroup, consumerHand.UpdateData)

	select {}
}
