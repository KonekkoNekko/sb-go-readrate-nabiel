package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"sb-go-readrate-nabiel/structs"
)

const categoryColumns = "id, name, created_at, created_by, modified_at, modified_by"

func RetrieveCategory(db *sql.DB, categoryID int) (category structs.Category, problem structs.Error) {
	queryStatement := fmt.Sprintf("SELECT %s FROM categories WHERE id = $1", categoryColumns)

	errQuery := db.QueryRow(queryStatement, categoryID).Scan(
		&category.ID, &category.Name, &category.CreatedAt, &category.CreatedBy, &category.ModifiedAt, &category.ModifiedBy)

	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return category, structs.Error{
				Message: fmt.Sprintf("category with identifier %d not found", categoryID),
				Status:  http.StatusNotFound,
			}
		}
		return category, structs.Error{
			Message: fmt.Sprintf("failed to retrieve category: %s", errQuery.Error()),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}

func RetrieveCategories(db *sql.DB) (categories []structs.Category, problem structs.Error) {
	queryStatement := fmt.Sprintf("SELECT %s FROM categories", categoryColumns)

	rowsResult, errQuery := db.Query(queryStatement)
	if errQuery != nil {
		return []structs.Category{}, structs.Error{
			Message: errQuery.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rowsResult.Close()

	for rowsResult.Next() {
		var categoryItem = structs.Category{}
		errScan := rowsResult.Scan(
			&categoryItem.ID, &categoryItem.Name, &categoryItem.CreatedAt, &categoryItem.CreatedBy, &categoryItem.ModifiedAt, &categoryItem.ModifiedBy)

		if errScan != nil {
			return []structs.Category{}, structs.Error{
				Message: errScan.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		categories = append(categories, categoryItem)
	}
	return
}

func StoreCategory(db *sql.DB, categoryData structs.Category) (problem structs.Error) {
	queryCommand := `INSERT INTO categories (name, created_by, modified_by) VALUES ($1, $2, $3)`

	_, errExec := db.Exec(queryCommand, categoryData.Name, categoryData.CreatedBy, categoryData.ModifiedBy)

	if errExec != nil {
		return structs.Error{
			Message: errExec.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}

func UpdateCategory(db *sql.DB, categoryData structs.Category) (problem structs.Error) {
	_, errExistence := RetrieveCategory(db, categoryData.ID)
	if errExistence.Message != "" {
		return errExistence
	}

	queryUpdate := `UPDATE categories SET name = $1, modified_by = $2, modified_at = NOW() WHERE id = $3`

	_, errExec := db.Exec(queryUpdate, categoryData.Name, categoryData.ModifiedBy, categoryData.ID)

	if errExec != nil {
		return structs.Error{
			Message: errExec.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}

func EraseCategory(db *sql.DB, categoryID int) (problem structs.Error) {
	_, errExistence := RetrieveCategory(db, categoryID)
	if errExistence.Message != "" {
		return errExistence
	}

	queryDelete := "DELETE FROM categories WHERE id = $1"
	_, errExec := db.Exec(queryDelete, categoryID)

	if errExec != nil {
		return structs.Error{
			Message: errExec.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return
}