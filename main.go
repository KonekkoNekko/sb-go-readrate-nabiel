package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sb-go-readrate-nabiel/controllers"
	"sb-go-readrate-nabiel/database"
	"sb-go-readrate-nabiel/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load environment variables
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Database connection string
	dbHost := os.Getenv("PGHOST")
	dbPort := os.Getenv("PGPORT")
	dbUser := os.Getenv("PGUSER")
	dbPassword := os.Getenv("PGPASSWORD")
	dbName := os.Getenv("PGDATABASE")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %s", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("Error closing database connection: %s", cerr)
		}
	}()

	// Ping database to verify connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}
	fmt.Println("Successfully connected to the database!")

	// Run database migrations
	database.DBMigrate(db)

	// Initialize Gin router
	router := gin.Default()

	// Public routes
	router.POST("/register", controllers.HandleRegisterUser) // Register new user
	router.GET("/books", controllers.HandleFindBooks)
	router.GET("/books/:id", controllers.HandleFindBook)
	router.GET("/categories", controllers.HandleFindCategories)
	router.GET("/categories/:id", controllers.HandleFindCategory)
	router.GET("/categories/:id/books", controllers.HandleFindBooksByCategory)
	router.GET("/books/:id/reviews", controllers.HandleGetReviewsForBook) // Get reviews for a specific book

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(middleware.Authenticate(db)) // Apply authentication middleware

	{
		// Book routes
		authenticated.POST("/books", controllers.HandleCreateBook)
		authenticated.PUT("/books/:id", controllers.HandleUpdateBook)
		authenticated.DELETE("/books/:id", controllers.HandleDeleteBook)

		// Category routes
		authenticated.POST("/categories", controllers.HandleCreateCategory)
		authenticated.PUT("/categories/:id", controllers.HandleUpdateCategory)
		authenticated.DELETE("/categories/:id", controllers.HandleDeleteCategory)

		// Review routes (NEW)
		authenticated.POST("/books/:book_id/reviews", controllers.HandleCreateReview) // Create review for a book
		authenticated.PUT("/reviews/:id", controllers.HandleUpdateReview)           // Update own review
		authenticated.DELETE("/reviews/:id", controllers.HandleDeleteReview)         // Delete own review
		authenticated.GET("/users/:user_id/reviews", controllers.HandleGetReviewsByUser) // Get reviews by a specific user (can be secured further if needed)
	}

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %s", err)
	}
}