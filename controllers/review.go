package controllers

import (
	"net/http"
	"sb-go-readrate-nabiel/database"
	"sb-go-readrate-nabiel/repository"
	"sb-go-readrate-nabiel/structs"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HandleCreateReview allows an authenticated user to create a review for a specific book.
func HandleCreateReview(context *gin.Context) {
	responseCode := http.StatusOK
	currentUser, exists := context.Get("user")
	if !exists {
		responseCode = http.StatusUnauthorized
		context.JSON(responseCode, gin.H{
			"detail": "Authentication required",
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

	bookID, err := strconv.Atoi(context.Param("book_id"))
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid book ID",
		})
		return
	}

	var newReview structs.Review
	if bindError := context.ShouldBindJSON(&newReview); bindError != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": bindError.Error(),
		})
		return
	}

	// Basic validation for rating
	if newReview.Rating < 1 || newReview.Rating > 5 {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Rating must be between 1 and 5",
		})
		return
	}

	newReview.BookID = bookID
	newReview.UserID = userRecord.ID
	newReview.CreatedBy = userRecord.ID
	newReview.ModifiedBy = userRecord.ID

	creationError := repository.StoreReview(database.DbConnection, newReview)

	if creationError.Message != "" {
		responseCode = creationError.Status
		context.JSON(responseCode, gin.H{
			"detail": creationError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"status":  "success",
		"message": "Review created successfully",
	})
}

// HandleGetReviewsForBook retrieves all reviews for a specified book.
func HandleGetReviewsForBook(context *gin.Context) {
	responseCode := http.StatusOK
	bookID, err := strconv.Atoi(context.Param("id")) // Note: param is "id" in main.go route for this
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid book ID",
		})
		return
	}

	reviewList, retrievalError := repository.RetrieveReviewsForBook(database.DbConnection, bookID)

	if retrievalError.Message != "" {
		responseCode = retrievalError.Status
		context.JSON(responseCode, gin.H{
			"detail": retrievalError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"items": reviewList,
	})
}

// HandleGetReviewsByUser retrieves all reviews written by a specified user.
func HandleGetReviewsByUser(context *gin.Context) {
	responseCode := http.StatusOK
	userID, err := strconv.Atoi(context.Param("user_id"))
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid user ID",
		})
		return
	}

	// Optional: Add authorization here if only authenticated user can see their own reviews,
	// or if only admins can see reviews of any user.
	// For now, assuming anyone can see reviews of any user by ID if they have the ID.

	reviewList, retrievalError := repository.RetrieveReviewsByUser(database.DbConnection, userID)

	if retrievalError.Message != "" {
		responseCode = retrievalError.Status
		context.JSON(responseCode, gin.H{
			"detail": retrievalError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"items": reviewList,
	})
}

// HandleUpdateReview allows an authenticated user to update their own review.
func HandleUpdateReview(context *gin.Context) {
	responseCode := http.StatusOK
	reviewID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid review ID",
		})
		return
	}

	currentUser, exists := context.Get("user")
	if !exists {
		responseCode = http.StatusUnauthorized
		context.JSON(responseCode, gin.H{
			"detail": "Authentication required",
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

	var updatedReview structs.Review
	if bindError := context.ShouldBindJSON(&updatedReview); bindError != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": bindError.Error(),
		})
		return
	}

	// Basic validation for rating
	if updatedReview.Rating < 1 || updatedReview.Rating > 5 {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Rating must be between 1 and 5",
		})
		return
	}

	// Retrieve the existing review to check ownership
	existingReview, retrieveErr := repository.FindSingleReview(database.DbConnection, reviewID)
	if retrieveErr.Message != "" {
		responseCode = retrieveErr.Status
		context.JSON(responseCode, gin.H{
			"detail": retrieveErr.Message,
		})
		return
	}

	// Authorization: Check if the current user is the owner of the review
	if existingReview.UserID != userRecord.ID {
		responseCode = http.StatusForbidden // 403 Forbidden
		context.JSON(responseCode, gin.H{
			"detail": "You are not authorized to update this review",
		})
		return
	}

	updatedReview.ID = reviewID
	updatedReview.ModifiedBy = userRecord.ID

	updateError := repository.UpdateExistingReview(database.DbConnection, updatedReview)

	if updateError.Message != "" {
		responseCode = updateError.Status
		context.JSON(responseCode, gin.H{
			"detail": updateError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"status":  "success",
		"message": "Review updated successfully",
	})
}

// HandleDeleteReview allows an authenticated user to delete their own review.
func HandleDeleteReview(context *gin.Context) {
	responseCode := http.StatusOK
	reviewID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid review ID",
		})
		return
	}

	currentUser, exists := context.Get("user")
	if !exists {
		responseCode = http.StatusUnauthorized
		context.JSON(responseCode, gin.H{
			"detail": "Authentication required",
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

	// Retrieve the existing review to check ownership
	existingReview, retrieveErr := repository.FindSingleReview(database.DbConnection, reviewID)
	if retrieveErr.Message != "" {
		responseCode = retrieveErr.Status
		context.JSON(responseCode, gin.H{
			"detail": retrieveErr.Message,
		})
		return
	}

	// Authorization: Check if the current user is the owner of the review
	if existingReview.UserID != userRecord.ID {
		responseCode = http.StatusForbidden // 403 Forbidden
		context.JSON(responseCode, gin.H{
			"detail": "You are not authorized to delete this review",
		})
		return
	}

	deletionError := repository.EraseReview(database.DbConnection, reviewID)

	if deletionError.Message != "" {
		responseCode = deletionError.Status
		context.JSON(responseCode, gin.H{
			"detail": deletionError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"status":  "success",
		"message": "Review deleted successfully",
	})
}