package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"sb-go-readrate-nabiel/structs"
)

const reviewColumns = "id, book_id, user_id, rating, comment, created_at, created_by, modified_at, modified_by"

// StoreReview saves a new review to the database.
func StoreReview(db *sql.DB, reviewData structs.Review) (problem structs.Error) {
	queryCommand := `INSERT INTO reviews (book_id, user_id, rating, comment, created_by, modified_by)
                     VALUES ($1, $2, $3, $4, $5, $6)`

	_, errExec := db.Exec(queryCommand, reviewData.BookID, reviewData.UserID, reviewData.Rating, reviewData.Comment, reviewData.CreatedBy, reviewData.ModifiedBy)

	if errExec != nil {
		// Check for unique constraint violation (a user can only review a book once)
		if errExec.Error() == `pq: duplicate key value violates unique constraint "reviews_book_id_user_id_key"` {
			return structs.Error{
				Message: fmt.Sprintf("user %d has already reviewed book %d", reviewData.UserID, reviewData.BookID),
				Status:  http.StatusConflict, // 409 Conflict
			}
		}
		return structs.Error{
			Message: errExec.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}

// RetrieveReviewsForBook fetches all reviews for a specific book.
func RetrieveReviewsForBook(db *sql.DB, bookID int) (collection []structs.Review, problem structs.Error) {
	queryStatement := fmt.Sprintf("SELECT %s FROM reviews WHERE book_id = $1 ORDER BY created_at DESC", reviewColumns)

	rowsData, errQuery := db.Query(queryStatement, bookID)
	if errQuery != nil {
		return []structs.Review{}, structs.Error{
			Message: errQuery.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rowsData.Close()

	for rowsData.Next() {
		var reviewItem = structs.Review{}
		errScan := rowsData.Scan(
			&reviewItem.ID, &reviewItem.BookID, &reviewItem.UserID, &reviewItem.Rating, &reviewItem.Comment,
			&reviewItem.CreatedAt, &reviewItem.CreatedBy, &reviewItem.ModifiedAt, &reviewItem.ModifiedBy,
		)

		if errScan != nil {
			return []structs.Review{}, structs.Error{
				Message: errScan.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		collection = append(collection, reviewItem)
	}

	if len(collection) == 0 {
		return collection, structs.Error{
			Message: fmt.Sprintf("no reviews found for book ID %d", bookID),
			Status:  http.StatusNotFound,
		}
	}
	return
}

// RetrieveReviewsByUser fetches all reviews written by a specific user.
func RetrieveReviewsByUser(db *sql.DB, userID int) (collection []structs.Review, problem structs.Error) {
	queryStatement := fmt.Sprintf("SELECT %s FROM reviews WHERE user_id = $1 ORDER BY created_at DESC", reviewColumns)

	rowsData, errQuery := db.Query(queryStatement, userID)
	if errQuery != nil {
		return []structs.Review{}, structs.Error{
			Message: errQuery.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rowsData.Close()

	for rowsData.Next() {
		var reviewItem = structs.Review{}
		errScan := rowsData.Scan(
			&reviewItem.ID, &reviewItem.BookID, &reviewItem.UserID, &reviewItem.Rating, &reviewItem.Comment,
			&reviewItem.CreatedAt, &reviewItem.CreatedBy, &reviewItem.ModifiedAt, &reviewItem.ModifiedBy,
		)

		if errScan != nil {
			return []structs.Review{}, structs.Error{
				Message: errScan.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		collection = append(collection, reviewItem)
	}

	if len(collection) == 0 {
		return collection, structs.Error{
			Message: fmt.Sprintf("no reviews found by user ID %d", userID),
			Status:  http.StatusNotFound,
		}
	}
	return
}

// FindSingleReview retrieves a single review by its ID.
func FindSingleReview(db *sql.DB, reviewID int) (record structs.Review, problem structs.Error) {
	querySingle := fmt.Sprintf("SELECT %s FROM reviews WHERE id = $1", reviewColumns)

	rowResult := db.QueryRow(querySingle, reviewID)

	errRow := rowResult.Scan(
		&record.ID, &record.BookID, &record.UserID, &record.Rating, &record.Comment,
		&record.CreatedAt, &record.CreatedBy, &record.ModifiedAt, &record.ModifiedBy,
	)

	if errRow != nil {
		if errRow == sql.ErrNoRows {
			return record, structs.Error{
				Message: fmt.Sprintf("review with identifier %d not found", reviewID),
				Status:  http.StatusNotFound,
			}
		}
		return record, structs.Error{
			Message: fmt.Sprintf("failed to retrieve review: %s", errRow.Error()),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}

// UpdateExistingReview updates an existing review.
func UpdateExistingReview(db *sql.DB, reviewData structs.Review) (problem structs.Error) {
	// Optional: Check if the review exists before updating
	_, errExistence := FindSingleReview(db, reviewData.ID)
	if errExistence.Message != "" {
		return errExistence
	}

	queryUpdate := `UPDATE reviews SET rating = $1, comment = $2, modified_by = $3, modified_at = NOW() WHERE id = $4`

	_, errExec := db.Exec(queryUpdate, reviewData.Rating, reviewData.Comment, reviewData.ModifiedBy, reviewData.ID)

	if errExec != nil {
		return structs.Error{
			Message: errExec.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}

// EraseReview deletes a review by its ID.
func EraseReview(db *sql.DB, reviewID int) (problem structs.Error) {
	// Optional: Check if the review exists before deleting
	_, errExistence := FindSingleReview(db, reviewID)
	if errExistence.Message != "" {
		return errExistence
	}

	queryDelete := "DELETE FROM reviews WHERE id = $1"
	_, errExec := db.Exec(queryDelete, reviewID)

	if errExec != nil {
		return structs.Error{
			Message: errExec.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}