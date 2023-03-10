package main

import (
	"quizard/constants"
	"quizard/database"
	"quizard/handlers"
	"quizard/routes"

	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "golang.org/x/text"
)

var (
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	constant := constants.New()

	// Parse command-line flags
	flag.Parse()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	dbConfig := database.Config{
		Host:     constant.DbHost,
		Port:     constant.DbPort,
		Password: constant.DbPassword,
		User:     constant.DbUser,
		DBName:   constant.DbName,
	}

	database.Connect(&dbConfig)

	database.Migrate(database.DB)

	// Bind routes
	routes.Routes(app, database.DB)

	// Setup static files
	app.Static("/", "./static/public")

	// Handle not founds
	app.Use(handlers.NotFound)

	// Listen on port set in .env
	log.Fatal(app.Listen(constant.Port))
}
