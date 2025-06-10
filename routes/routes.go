package routes

import (
	"github.com/gofiber/fiber/v2"
	"task-manager/handlers"
)

func SetupRoutes(app *fiber.App) {
	//api group
	api := app.Group("/api")

	//task routes
	tasks := api.Group("/tasks")
	tasks.Get("/", handlers.GetTask)
	tasks.Get("/:id", handlers.GetTask)

}
