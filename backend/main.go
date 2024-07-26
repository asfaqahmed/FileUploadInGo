package main

import (
	"fmt"
	"log"

	// "os"
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
	dsn := "host=localhost user=postgress dbname=file_upload password=leeladealwis@1111"
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

	log.Fatal(app.Listen(":3000"))
}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
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
	})
}
