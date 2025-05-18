package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"sb-go-readrate-nabiel/structs"

	"golang.org/x/crypto/bcrypt"
)

// StoreUser saves a new user to the database after hashing their password.
func StoreUser(db *sql.DB, userData structs.User) (problem structs.Error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return structs.Error{
			Message: fmt.Sprintf("failed to hash password: %s", err.Error()),
			Status:  http.StatusInternalServerError,
		}
	}

	queryCommand := `INSERT INTO users (username, password, created_by, modified_by) VALUES ($1, $2, $3, $4)`

	// For initial user creation, created_by/modified_by can be 0 or the user's own ID if known.
	// For simplicity, we set it to 0 assuming system creation or self-creation.
	_, errExec := db.Exec(queryCommand, userData.Username, string(hashedPassword), 0, 0)

	if errExec != nil {
		// Check for unique constraint violation (e.g., duplicate username)
		if errExec.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return structs.Error{
				Message: "username already exists",
				Status:  http.StatusConflict, // 409 Conflict
			}
		}
		return structs.Error{
			Message: errExec.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return structs.Error{} // No error
}

// FindUserByID retrieves a single user by their ID. (Useful for getting user info after auth)
func FindUserByID(db *sql.DB, userID int) (record structs.User, problem structs.Error) {
	querySingle := "SELECT id, username, created_at, created_by, modified_at, modified_by FROM users WHERE id = $1"

	rowResult := db.QueryRow(querySingle, userID)

	errRow := rowResult.Scan(
		&record.ID,
		&record.Username,
		&record.CreatedAt,
		&record.CreatedBy,
		&record.ModifiedAt,
		&record.ModifiedBy,
	)

	if errRow != nil {
		if errRow == sql.ErrNoRows {
			return record, structs.Error{
				Message: fmt.Sprintf("user with identifier %d not found", userID),
				Status:  http.StatusNotFound,
			}
		}
		return record, structs.Error{
			Message: fmt.Sprintf("failed to retrieve user: %s", errRow.Error()),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}