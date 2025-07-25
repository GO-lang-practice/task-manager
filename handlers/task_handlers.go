package handlers

import (
	"context"
	"net/http"
	"task-manager/database"
	"task-manager/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var taskCollection *mongo.Collection

func InitHandlers() {
	taskCollection = database.GetCollection("tasks")
}

// get all tasks
func GetTasks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := taskCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tasks",
		})
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode tasks",
		})
	}

	return c.JSON(fiber.Map{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// get a single task by id
func GetTask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	var task models.Task
	err = taskCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch task",
		})
	}

	return c.JSON(task)
}

// create a new task
func CreateTask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var request models.CreateTaskRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if request.Title == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	task := models.Task{
		Id:          primitive.NewObjectID(),
		Title:       request.Title,
		Description: request.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := taskCollection.InsertOne(ctx, task)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create task",
		})
	}

	return c.Status(http.StatusCreated).JSON(task)
}

// UpdateTask - Update an existing task
func UpdateTask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	var request models.UpdateTaskRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	if request.Title != "" {
		update["$set"].(bson.M)["title"] = request.Title
	}
	if request.Description != "" {
		update["$set"].(bson.M)["description"] = request.Description
	}
	if request.Completed != nil {
		update["$set"].(bson.M)["completed"] = *request.Completed
	}

	result, err := taskCollection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update task",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	// Fetch and return updated task
	var updatedTask models.Task
	err = taskCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&updatedTask)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch updated task",
		})
	}

	return c.JSON(updatedTask)
}
