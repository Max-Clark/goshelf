package postgresql

import (
	"errors"
	"fmt"
	"strings"
	"time"

	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
	_ "github.com/lib/pq"
)

// Creates a new book in the database. Returns the book_id generated.
func (pg *PgDb) BookCreate(b *v1.Book) (*int, error) {
	if b == nil {
		return nil, nil
	}

	authId, err := pg.CreateAuthorIfNew(b)

	if err != nil {
		return nil, err
	}

	// TODO: automate this with tags
	inserts := []string{"title", "author_id"}
	queryValues := []interface{}{
		b.Title,
		authId,
	}

	if b.PublishDate != nil {
		inserts = append(inserts, "publish_date")
		queryValues = append(queryValues, b.PublishDate.Format(time.RFC3339))
	}

	if b.Edition != nil {
		inserts = append(inserts, "edition")
		queryValues = append(queryValues, b.Edition)
	}

	if b.Description != nil {
		inserts = append(inserts, "description")
		queryValues = append(queryValues, b.Description)
	}

	if b.Genre != nil {
		inserts = append(inserts, "genre")
		queryValues = append(queryValues, b.Genre)
	}

	valueVars := make([]string, len(inserts))
	for i := 0; i < len(valueVars); i++ {
		valueVars[i] = "$" + fmt.Sprint(i+1)
	}

	queryStr := fmt.Sprintf(`
		INSERT INTO %s.book ( %s )
		VALUES              ( %s )
		RETURNING book_id
	`,
		pg.SchemaVersion,
		strings.Join(inserts, " , "),
		strings.Join(valueVars, " , "),
	)

	rows, err := pg.SqlDb.Query(
		queryStr,
		queryValues...,
	)

	if err != nil {
		return nil, err
	}

	return ScanReturnedId(rows)
}

// Returns a book from the database based on id.
func (pg *PgDb) BookGet(id int) (*v1.Book, error) {

	queryStr := fmt.Sprintf(`
		SELECT b.book_id, b.created_ts, b.title, b.publish_date, b.edition, b.description, b.genre, a.author_id, a.created_ts, a.first_name, a.last_name
		FROM %s.book b 
		INNER JOIN %s.author a ON b.author_id = a.author_id 
		WHERE b.book_id = $1
	`, pg.SchemaVersion, pg.SchemaVersion)

	rows, err := pg.SqlDb.Query(queryStr, id)

	if err != nil {
		return nil, err
	}

	books, err := ScanReturnedBooks(rows)

	if err != nil {
		return nil, err
	}

	if len(books) < 1 {
		return nil, nil
	}

	return &books[0], nil
}

// Removes a book from the database based on ID.
func (pg *PgDb) BookRemove(id int) error {
	// the SQL object doesn't return rows adjusted, so we'll check
	// to see if the book exists and error if not
	book, err := pg.BookGet(id)

	if err != nil {
		return err
	}

	if book == nil {
		return errors.New("book not found")
	}

	queryStr := fmt.Sprintf(`
		DELETE FROM %s.book b 
		WHERE b.book_id = $1
	`, pg.SchemaVersion)

	rows, err := pg.SqlDb.Query(queryStr, id)

	if err != nil {
		return err
	}

	// Close the connection once we're done
	defer rows.Close()

	return nil
}

// Returns an array of books based on filter. If no filters given,
// this function returns all books. Title and genre are wildcard searches,
// edition is equality.
func (pg *PgDb) BookFilter(title *string, genre *string, edition *int) ([]v1.Book, error) {
	queryStr := fmt.Sprintf(`
		SELECT b.book_id, b.created_ts, b.title, b.publish_date, b.edition, b.description, b.genre, a.author_id, a.created_ts, a.first_name, a.last_name
		FROM %s.book b 
		INNER JOIN %s.author a ON b.author_id = a.author_id 
	`, pg.SchemaVersion, pg.SchemaVersion)

	// TODO: Automate this section based on reflection/validation
	// The next section generates a dynamic where string
	// (e.g., where x = y and y = z)
	wheres := make([]string, 0)
	values := make([]interface{}, 0)
	idx := 1

	// Perform wildcard search on title if given
	if title != nil {
		wheres = append(wheres, " b.title LIKE '%' || $"+fmt.Sprint(idx)+" || '%' ")
		values = append(values, *title)
		idx++
	}

	// Perform wildcard search on genre if given
	if genre != nil {
		wheres = append(wheres, " b.genre LIKE '%' || $"+fmt.Sprint(idx)+" || '%' ")
		values = append(values, *genre)
		idx++
	}

	// Perform equality check on title if given
	if edition != nil {
		wheres = append(wheres, " b.edition = $"+fmt.Sprint(idx)+" ")
		values = append(values, *edition)
		idx++
	}

	if len(wheres) > 0 {
		queryStr += " WHERE " + strings.Join(wheres, " AND ")
	}

	rows, err := pg.SqlDb.Query(queryStr, values...)

	if err != nil {
		return nil, err
	}

	return ScanReturnedBooks(rows)
}
