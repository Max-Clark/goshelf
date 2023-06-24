package postgresql

import (
	"fmt"

	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
	_ "github.com/lib/pq"
)

// Gets an author by name (i.e., author.first_name & author.last_name full match).
// Returns nil, nil for no rows found or nil book pointer.
func (pg *PgDb) GetAuthorByName(b *v1.Book) (*v1.Author, error) {
	if b == nil {
		return nil, nil
	}

	queryStr := fmt.Sprintf(`
		SELECT a.author_id, a.created_ts, a.first_name, a.last_name FROM %s.author a WHERE first_name = $1 AND last_name = $2
	`, pg.SchemaVersion)

	rows, err := pg.SqlDb.Query(
		queryStr,
		b.Author.FirstName,
		b.Author.LastName,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if ok := rows.Next(); !ok {
		return nil, nil
	}

	a := v1.Author{}

	err = rows.Scan(
		&a.AuthorId,
		&a.CreatedTs,
		&a.FirstName,
		&a.LastName,
	)

	if err != nil {
		return nil, err
	}

	return &a, nil

}

// If an author does not exist (i.e., first_name and last_name found in database)
// create it. Returns the created or existing author_id.
func (pg *PgDb) CreateAuthorIfNew(b *v1.Book) (*int, error) {
	if b == nil {
		return nil, nil
	}

	auth, err := pg.GetAuthorByName(b)

	if err != nil {
		return nil, err
	}

	if auth != nil {
		return &auth.AuthorId, nil
	}

	queryStr := fmt.Sprintf(`
		INSERT INTO %s.author (first_name, last_name)
		VALUES ($1, $2)
		RETURNING author_id
	`, pg.SchemaVersion)

	rows, err := pg.SqlDb.Query(
		queryStr,
		b.Author.FirstName,
		b.Author.LastName,
	)

	if err != nil {
		return nil, err
	}

	return ScanReturnedId(rows)

}
