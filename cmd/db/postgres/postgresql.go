package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	db "github.com/Max-Clark/goshelf/cmd/db"
	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
	_ "github.com/lib/pq"
)

type PgDb struct {
	SqlDb         *sql.DB
	SchemaVersion string // Used for migrations
	Config        db.ConnectionConfig
}

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

// Creates a new book in the database. Returns the book_id generated.
func (pg *PgDb) BookCreate(b *v1.Book) (*int, error) {
	if b == nil {
		return nil, nil
	}

	authId, err := pg.CreateAuthorIfNew(b)

	if err != nil {
		return nil, err
	}

	queryStr := fmt.Sprintf(`
		INSERT INTO %s.book (title,publish_date,edition,description,genre,author_id)
		VALUES              ($1,   $2,          $3,     $4,         $5,   $6       )
		RETURNING book_id
	`, pg.SchemaVersion)

	publishDate := b.PublishDate.Format(time.RFC3339)

	rows, err := pg.SqlDb.Query(
		queryStr,
		b.Title,
		publishDate,
		*b.Edition,
		*b.Description,
		*b.Genre,
		authId,
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
func (pg *PgDb) CollectionCreate(title *string, bookIds *[]int) (*int, error) {
	return nil, nil
}

func (pg *PgDb) CollectionGet(title string) (*v1.Collection, error) {
	queryStr := fmt.Sprintf(`
		SELECT * FROM %s.collection c WHERE c.title = $1
	`, pg.SchemaVersion)

	rows, err := pg.SqlDb.Query(queryStr, title)

	if err != nil {
		return nil, err
	}

	collections, err := ScanReturnedCollections(rows)

	if err != nil {
		return nil, err
	}

	if len(collections) < 1 {
		return nil, nil
	}

	collection := &collections[0]

	queryStr = fmt.Sprintf(`
		SELECT 	b.book_id, b.created_ts, b.title, b.publish_date, b.edition, b.description, 
				b.genre, a.author_id, a.created_ts, a.first_name, a.last_name
		FROM %s.collection_books cb
		INNER JOIN %s.book b ON cb.book_id = b.book_id
		INNER JOIN %s.author a ON b.author_id = a.author_id 
	`, pg.SchemaVersion, pg.SchemaVersion, pg.SchemaVersion)

	rows, err = pg.SqlDb.Query(queryStr, title)

	if err != nil {
		return nil, err
	}

	books, err := ScanReturnedBooks(rows)

	if err != nil {
		return nil, err
	}

	collection.Books = books

	return collection, nil
}

func (pg *PgDb) CollectionRemove(title string) error {
	queryStr := fmt.Sprintf(`
		DELETE FROM %s.collection c
		WHERE c.title = $1
	`, pg.SchemaVersion)

	rows, err := pg.SqlDb.Query(queryStr, title)

	if err != nil {
		return err
	}

	// Close the connection once we're done
	defer rows.Close()

	return nil
}

func (pg *PgDb) Connect() error {
	if pg.SqlDb != nil {
		return nil
	}

	host := pg.Config.Host
	port := pg.Config.Port
	user := pg.Config.User
	password := pg.Config.Password
	dbname := pg.Config.DbName
	sslmode := pg.Config.SslMode

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connString)

	if err != nil {
		return err
	}

	pg.SqlDb = db

	return nil
}
