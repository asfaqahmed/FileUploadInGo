package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type File struct {
	ID        uint   `gorm:"primary_key"`
	Filename  string `gorm:"unique;not null"`
	CreatedAt time.Time
}

var db *gorm.DB

func initDatabase() {
	var err error
	dsn := "host=localhost user=postgres dbname=file_upload sslmode=disable password=leeladealwis@1111"

	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	fmt.Println("Database connection successfully opened")

	db.AutoMigrate(&File{})
	fmt.Println("Database Migrated")
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	initDatabase()
	defer db.Close()

	app.Post("/upload", uploadFile)
	app.Get("/files/:id", getFile)

	log.Fatal(app.Listen(":8081"))
}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Ensure uploads directory exists
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", os.ModePerm)
	}

	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return err
	}

	newFile := File{
		Filename:  file.Filename,
		CreatedAt: time.Now(),
	}
	db.Create(&newFile)

	return c.JSON(fiber.Map{
		"message": "File uploaded successfully",
		"fileID":  newFile.ID,
	})
}

func getFile(c *fiber.Ctx) error {
	// Retrieve the file ID from the URL parameters
	id := c.Params("id")

	// Retrieve the file record from the database
	var file File
	if err := db.First(&file, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Construct the file path
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)

	// Serve the file
	return c.SendFile(filePath)
}
