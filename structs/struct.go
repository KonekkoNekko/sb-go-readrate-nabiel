package structs

import "time"

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ReleaseYear int       `json:"release_year"`
	Price       int       `json:"price"`
	TotalPage   int       `json:"total_page"`
	Thickness   string    `json:"thickness"`
	CategoryID  int       `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   int       `json:"created_by"` // Changed to int
	ModifiedAt  time.Time `json:"modified_at"`
	ModifiedBy  int       `json:"modified_by"` // Changed to int
}

type Category struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  int       `json:"created_by"` // Changed to int
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy int       `json:"modified_by"` // Changed to int
}

type User struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"` // This will store the hashed password
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  int       `json:"created_by"` // Changed to int (refers to itself for initial creation, or 0/null)
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy int       `json:"modified_by"` // Changed to int (refers to itself for initial modification, or 0/null)
}

// NEW Struct: Review
type Review struct {
	ID         int       `json:"id"`
	BookID     int       `json:"book_id"`
	UserID     int       `json:"user_id"`
	Rating     int       `json:"rating"` // e.g., 1-5 stars
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  int       `json:"created_by"` // Changed to int
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy int       `json:"modified_by"` // Changed to int
}

type Error struct {
	Message string `json:"detail"`
	Status  int    `json:"status"`
}