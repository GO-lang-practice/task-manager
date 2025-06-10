package routes

import (
	"task-manager/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	//api group
	api := app.Group("/api")

	tasks := api.Group("/tasks")
	tasks.Get("/", handlers.GetTasks)    // Get all tasks
	tasks.Get("/:id", handlers.GetTask)  // Get single task by ID
	tasks.Post("/", handlers.CreateTask) // Create new task

}
