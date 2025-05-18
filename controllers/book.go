package controllers

import (
	"net/http"
	"sb-go-readrate-nabiel/database"
	"sb-go-readrate-nabiel/repository"
	"sb-go-readrate-nabiel/structs"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleFindBooks(context *gin.Context) {
	responseCode := http.StatusOK
	bookList, retrievalError := repository.RetrieveBooks(database.DbConnection)

	if retrievalError.Message != "" {
		responseCode = retrievalError.Status
		context.JSON(responseCode, gin.H{
			"detail": retrievalError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"items": bookList,
	})
}

func HandleFindBook(context *gin.Context) {
	responseCode := http.StatusOK
	bookID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid book ID",
		})
		return
	}

	singleBook, retrievalError := repository.FindSingleBook(database.DbConnection, bookID)

	if retrievalError.Message != "" {
		responseCode = retrievalError.Status
		context.JSON(responseCode, gin.H{
			"detail": retrievalError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"item": singleBook,
	})
}

func HandleCreateBook(context *gin.Context) {
	responseCode := http.StatusOK
	currentUser, exists := context.Get("user")
	if !exists {
		responseCode = http.StatusUnauthorized
		context.JSON(responseCode, gin.H{
			"detail": "User not authenticated",
		})
		return
	}
	userRecord, ok := currentUser.(structs.User)
	if !ok {
		responseCode = http.StatusInternalServerError
		context.JSON(responseCode, gin.H{
			"detail": "Failed to get user data from context",
		})
		return
	}

	var newBook structs.Book
	if bindError := context.ShouldBindJSON(&newBook); bindError != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": bindError.Error(),
		})
		return
	}

	if newBook.ReleaseYear < 1980 || newBook.ReleaseYear > 2025 { // Adjusted to 2025 for current year + future
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Release year must be within the range 1980-2025",
		})
		return
	}

	if newBook.TotalPage > 100 {
		newBook.Thickness = "tebal"
	} else {
		newBook.Thickness = "tipis"
	}

	newBook.CreatedBy = userRecord.ID // Changed to int
	newBook.ModifiedBy = userRecord.ID // Changed to int

	creationError := repository.StoreBook(database.DbConnection, newBook)

	if creationError.Message != "" {
		responseCode = creationError.Status
		context.JSON(responseCode, gin.H{
			"detail": creationError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"status": "success",
		"message": "Book created successfully",
	})
}

func HandleUpdateBook(context *gin.Context) {
	responseCode := http.StatusOK
	bookID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid book ID",
		})
		return
	}

	currentUser, exists := context.Get("user")
	if !exists {
		responseCode = http.StatusUnauthorized
		context.JSON(responseCode, gin.H{
			"detail": "User not authenticated",
		})
		return
	}
	userRecord, ok := currentUser.(structs.User)
	if !ok {
		responseCode = http.StatusInternalServerError
		context.JSON(responseCode, gin.H{
			"detail": "Failed to get user data from context",
		})
		return
	}

	var updatedBook structs.Book
	if bindError := context.ShouldBindJSON(&updatedBook); bindError != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": bindError.Error(),
		})
		return
	}

	if updatedBook.ReleaseYear < 1980 || updatedBook.ReleaseYear > 2025 { // Adjusted to 2025
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Release year must be between 1980 and 2025",
		})
		return
	}

	if updatedBook.TotalPage > 100 {
		updatedBook.Thickness = "tebal"
	} else {
		updatedBook.Thickness = "tipis"
	}

	updatedBook.ID = bookID
	updatedBook.ModifiedBy = userRecord.ID // Changed to int

	updateError := repository.UpdateExistingBook(database.DbConnection, updatedBook)

	if updateError.Message != "" {
		responseCode = updateError.Status
		context.JSON(responseCode, gin.H{
			"detail": updateError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"status": "success",
		"message": "Book updated successfully",
	})
}

func HandleDeleteBook(context *gin.Context) {
	responseCode := http.StatusOK
	bookID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid book ID",
		})
		return
	}

	deletionError := repository.EraseBook(database.DbConnection, bookID)

	if deletionError.Message != "" {
		responseCode = deletionError.Status
		context.JSON(responseCode, gin.H{
			"detail": deletionError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"status": "success",
		"message": "Book deleted successfully",
	})
}