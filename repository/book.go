package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"sb-go-readrate-nabiel/structs"
)

const bookColumns = "id, title, category_id, description, image_url, release_year, price, total_page, thickness, created_at, created_by, modified_at, modified_by"

func RetrieveBooks(db *sql.DB) (collection []structs.Book, problem structs.Error) {
	queryStatement := fmt.Sprintf("SELECT %s FROM books", bookColumns)

	rowsData, errQuery := db.Query(queryStatement)
	if errQuery != nil {
		return []structs.Book{}, structs.Error{
			Message: errQuery.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rowsData.Close()

	for rowsData.Next() {
		var bookItem = structs.Book{}
		errScan := rowsData.Scan(
			&bookItem.ID,
			&bookItem.Title,
			&bookItem.CategoryID,
			&bookItem.Description,
			&bookItem.ImageURL,
			&bookItem.ReleaseYear,
			&bookItem.Price,
			&bookItem.TotalPage,
			&bookItem.Thickness,
			&bookItem.CreatedAt,
			&bookItem.CreatedBy,
			&bookItem.ModifiedAt,
			&bookItem.ModifiedBy,
		)

		if errScan != nil {
			return []structs.Book{}, structs.Error{
				Message: errScan.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		collection = append(collection, bookItem)
	}
	return
}

func StoreBook(db *sql.DB, bookData structs.Book) (problem structs.Error) {
	queryCommand := `INSERT INTO books (title, category_id, description, image_url, release_year, price, total_page, thickness, created_by, modified_by)
                     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, errExec := db.Exec(queryCommand, bookData.Title, bookData.CategoryID, bookData.Description, bookData.ImageURL, bookData.ReleaseYear, bookData.Price, bookData.TotalPage, bookData.Thickness, bookData.CreatedBy, bookData.ModifiedBy)

	if errExec != nil {
		return structs.Error{
			Message: errExec.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}

func FindSingleBook(db *sql.DB, bookID int) (record structs.Book, problem structs.Error) {
	querySingle := fmt.Sprintf("SELECT %s FROM books WHERE id = $1", bookColumns)

	rowResult := db.QueryRow(querySingle, bookID)

	errRow := rowResult.Scan(
		&record.ID,
		&record.Title,
		&record.CategoryID,
		&record.Description,
		&record.ImageURL,
		&record.ReleaseYear,
		&record.Price,
		&record.TotalPage,
		&record.Thickness,
		&record.CreatedAt,
		&record.CreatedBy,
		&record.ModifiedAt,
		&record.ModifiedBy,
	)

	if errRow != nil {
		if errRow == sql.ErrNoRows {
			return record, structs.Error{
				Message: fmt.Sprintf("book with identifier %d not found", bookID),
				Status:  http.StatusNotFound,
			}
		}
		return record, structs.Error{
			Message: fmt.Sprintf("failed to retrieve book: %s", errRow.Error()),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}

func UpdateExistingBook(db *sql.DB, bookData structs.Book) (problem structs.Error) {
	_, errExistence := FindSingleBook(db, bookData.ID)
	if errExistence.Message != "" {
		return errExistence
	}

	queryUpdate := `UPDATE books SET title = $1, category_id = $2, description = $3, image_url = $4,
                    release_year = $5, price = $6, total_page = $7, thickness = $8, modified_by = $9,
                    modified_at = NOW() WHERE id = $10`

	_, errExec := db.Exec(queryUpdate, bookData.Title, bookData.CategoryID, bookData.Description, bookData.ImageURL,
		bookData.ReleaseYear, bookData.Price, bookData.TotalPage, bookData.Thickness, bookData.ModifiedBy, bookData.ID)

	if errExec != nil {
		return structs.Error{
			Message: errExec.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}

func EraseBook(db *sql.DB, bookID int) (problem structs.Error) {
	_, errExistence := FindSingleBook(db, bookID)
	if errExistence.Message != "" {
		return errExistence
	}

	queryDelete := "DELETE FROM books WHERE id = $1"
	_, errExec := db.Exec(queryDelete, bookID)

	if errExec != nil {
		return structs.Error{
			Message: errExec.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}

func RetrieveBooksForCategory(db *sql.DB, categoryIdentifier int) (collection []structs.Book, problem structs.Error) {
	queryCategory := fmt.Sprintf("SELECT %s FROM books WHERE category_id = $1", bookColumns)

	rowsData, errQuery := db.Query(queryCategory, categoryIdentifier)
	if errQuery != nil {
		return collection, structs.Error{
			Message: fmt.Sprintf("failed to retrieve books: %s", errQuery.Error()),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rowsData.Close()

	itemCount := 0
	for rowsData.Next() {
		itemCount++
		var bookItem = structs.Book{}

		// Ensure scan order matches the selected columns
		errScan := rowsData.Scan(
			&bookItem.ID,
			&bookItem.Title,       // Changed order to match bookColumns
			&bookItem.CategoryID,
			&bookItem.Description,
			&bookItem.ImageURL,
			&bookItem.ReleaseYear,
			&bookItem.Price,
			&bookItem.TotalPage,
			&bookItem.Thickness,
			&bookItem.CreatedAt,
			&bookItem.CreatedBy,
			&bookItem.ModifiedAt,
			&bookItem.ModifiedBy,
		)

		if errScan != nil {
			return collection, structs.Error{
				Message: errScan.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		collection = append(collection, bookItem)
	}

	if itemCount == 0 {
		return collection, structs.Error{
			Message: fmt.Sprintf("no book found under Category ID %d", categoryIdentifier),
			Status:  http.StatusNotFound,
		}
	}
	return
}