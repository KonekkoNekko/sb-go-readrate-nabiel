package controllers

import (
	"net/http"
	"sb-go-readrate-nabiel/database"
	"sb-go-readrate-nabiel/repository"
	"sb-go-readrate-nabiel/structs"

	"github.com/gin-gonic/gin"
)

// HandleRegisterUser handles user registration.
func HandleRegisterUser(context *gin.Context) {
	responseCode := http.StatusOK
	var newUser structs.User
	if bindError := context.ShouldBindJSON(&newUser); bindError != nil {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": bindError.Error(),
		})
		return
	}

	// Basic validation
	if newUser.Username == "" || newUser.Password == "" {
		responseCode = http.StatusBadRequest
		context.JSON(responseCode, gin.H{
			"detail": "Username and password cannot be empty",
		})
		return
	}

	// Store user (password will be hashed inside repository)
	creationError := repository.StoreUser(database.DbConnection, newUser)

	if creationError.Message != "" {
		responseCode = creationError.Status
		context.JSON(responseCode, gin.H{
			"detail": creationError.Message,
		})
		return
	}

	context.JSON(responseCode, gin.H{
		"status":  "success",
		"message": "User registered successfully",
	})
}