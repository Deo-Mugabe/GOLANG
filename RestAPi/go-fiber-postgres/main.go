package main

import (
	"context"
	"log"
	"os"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"githumb.com/Deo-Mugabe/go-fiber-postgres/models"
	"githumb.com/Deo-Mugabe/go-fiber-postgres/storage"
	"golang.org/x/mod/sumdb/storage"
	"gorm.io/gorm" // helps to interract with the database
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error // creating a book using datbase
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has bee adde"})

	return nil
}

func (r *Repository) DeleteBook(context, *fiber.Ctx) error {
	bookModels := &[]models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(bookModels, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"messge": "could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book delete successfully",
	})
	return nil
}

func (r *Repository) Getbooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get books",
			"data":    bookModels,
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book fetched successfully",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) GetBookByID(context, *fiber.Ctx) error {

	id := context.Params("id")
	bookModel := &models.Books{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Prinln("the ID is", id)

	err := r.DB.Where("id =?", id).First(bookModel).Error
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"messge": "could not geet the book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book fetched successfully",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book", r.DeleteBook)
	api.Get("/get_book/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not establish connection")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate db")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
