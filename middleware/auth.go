package middleware

import (
	"database/sql"
	"net/http"
	"sb-go-readrate-nabiel/structs"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt" // Import bcrypt
)

func Authenticate(dbs *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		username, password, present := context.Request.BasicAuth()

		if !present {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "Authentication required",
			})
			return
		}

		var account structs.User
		// Select the hashed password
		statement := "SELECT id, username, password FROM users WHERE username = $1"
		err := dbs.QueryRow(statement, username).Scan(&account.ID, &account.Username, &account.Password)

		if err != nil {
			if err == sql.ErrNoRows {
				context.JSON(http.StatusUnauthorized, gin.H{"detail": "Incorrect credentials"})
			} else {
				context.JSON(http.StatusInternalServerError, gin.H{"detail": "Authentication failed: " + err.Error()})
			}
			context.Abort()
			return
		}

		// Compare the provided password with the hashed password from DB
		err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"detail": "Incorrect credentials"})
			context.Abort()
			return
		}

		// Authentication successful, set the user in context
		context.Set("user", account) // Changed "account" to "user" for consistency
		context.Next()
	}
}