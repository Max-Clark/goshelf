package postgresql

import (
	"database/sql"

	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
)

// Scans rows for one row expecting one integer parameter.
// Returns nil, nil for no rows returned or nil rows pointer.
func ScanReturnedCollections(rows *sql.Rows) ([]v1.Collection, error) {
	if rows == nil {
		return nil, nil
	}

	collections := make([]v1.Collection, 0)

	// Close the connection once we're done
	defer rows.Close()

	collection := &v1.Collection{}

	for rows.Next() {

		err := rows.Scan(
			&collection.Title,
			&collection.CreatedTs,
		)

		if err != nil {
			return nil, err
		}

		collections = append(collections, *collection)
	}

	return collections, nil
}

// Scans rows for one row expecting one integer parameter.
// Returns nil, nil for no rows returned or nil rows pointer.
func ScanReturnedId(rows *sql.Rows) (*int, error) {
	if rows == nil {
		return nil, nil
	}

	defer rows.Close()

	if ok := rows.Next(); !ok {
		return nil, nil
	}

	id := 0
	err := rows.Scan(
		&id,
	)

	if err != nil {
		return nil, err
	}

	return &id, nil
}

// Returns a series of books returned by rows. Returns an empty array
// if no rows returned.
func ScanReturnedBooks(rows *sql.Rows) ([]v1.Book, error) {
	if rows == nil {
		return nil, nil
	}

	books := make([]v1.Book, 0)

	// Close the connection once we're done
	defer rows.Close()

	for rows.Next() {

		book := &v1.Book{}
		err := rows.Scan(
			&book.BookId,
			&book.CreatedTs,
			&book.Title,
			&book.PublishDate,
			&book.Edition,
			&book.Description,
			&book.Genre,
			&book.Author.AuthorId,
			&book.Author.CreatedTs,
			&book.Author.FirstName,
			&book.Author.LastName,
		)

		if err != nil {
			return nil, err
		}

		books = append(books, *book)
	}

	return books, nil
}
