package controllers

import (
	"net/http"
	"sb-go-readrate-nabiel/database"
	"sb-go-readrate-nabiel/repository"
	"sb-go-readrate-nabiel/structs"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleFindCategories(context *gin.Context) {
	responseCode := http.StatusOK
	categoryList, retrievalError := repository.RetrieveCategories(database.DbConnection)

	if retrievalError.Message != "" {
		responseCode = retrievalError.Status
		context.JSON(responseCode, gin.H{
			"detail": retrievalError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"items": categoryList,
	})
}

func HandleCreateCategory(context *gin.Context) {
	responseCode := http.StatusOK
	var newCategory structs.Category
	if bindError := context.ShouldBindJSON(&newCategory); bindError != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": bindError.Error(),
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

	newCategory.CreatedBy = userRecord.ID   // Changed to int
	newCategory.ModifiedBy = userRecord.ID // Changed to int

	creationError := repository.StoreCategory(database.DbConnection, newCategory)

	if creationError.Message != "" {
		responseCode = creationError.Status
		context.JSON(responseCode, gin.H{
			"detail": creationError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"status": "success",
		"message": "Category created successfully",
	})
}

func HandleUpdateCategory(context *gin.Context) {
	responseCode := http.StatusOK
	categoryID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid category ID",
		})
		return
	}

	var updatedCategory structs.Category
	if bindError := context.ShouldBindJSON(&updatedCategory); bindError != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": bindError.Error(),
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

	updatedCategory.ID = categoryID
	// Fix: updatedCategory.CreatedBy should not be set here on update unless it's a specific logic.
	// We only need to set ModifiedBy for updates.
	// updatedCategory.CreatedBy = userRecord.ID
	updatedCategory.ModifiedBy = userRecord.ID // Changed to int

	updateError := repository.UpdateCategory(database.DbConnection, updatedCategory)

	if updateError.Message != "" {
		responseCode = updateError.Status
		context.JSON(responseCode, gin.H{
			"detail": updateError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"status": "success",
		"message": "Category updated successfully",
	})
}

func HandleDeleteCategory(context *gin.Context) {
	responseCode := http.StatusOK
	categoryID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid category ID",
		})
		return
	}

	deletionError := repository.EraseCategory(database.DbConnection, categoryID)

	if deletionError.Message != "" {
		responseCode = deletionError.Status
		context.JSON(responseCode, gin.H{
			"detail": deletionError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"status": "success",
		"message": "Category deleted successfully",
	})
}

func HandleFindCategory(context *gin.Context) {
	responseCode := http.StatusOK
	categoryID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid category ID",
		})
		return
	}

	singleCategory, retrievalError := repository.RetrieveCategory(database.DbConnection, categoryID)

	if retrievalError.Message != "" {
		responseCode = retrievalError.Status
		context.JSON(responseCode, gin.H{
			"detail": retrievalError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"item": singleCategory,
	})
}

func HandleFindBooksByCategory(context *gin.Context) {
	responseCode := http.StatusOK
	categoryID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Invalid category ID",
		})
		return
	}

	booksInCategory, retrievalError := repository.RetrieveBooksForCategory(database.DbConnection, categoryID)

	if retrievalError.Message != "" {
		responseCode = retrievalError.Status
		context.JSON(responseCode, gin.H{
			"detail": retrievalError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"items": booksInCategory,
	})
}