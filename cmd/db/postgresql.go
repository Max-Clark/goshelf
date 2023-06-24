package db

import (
	"database/sql"
	"fmt"
	"time"

	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
	_ "github.com/lib/pq"
)

type PgDb struct {
	SqlDb         *sql.DB
	SchemaVersion string // Used for migrations
	Config        ConnectionConfig
}

func ScanReturnedId(rows *sql.Rows) (*int, error) {
	id := 0

	defer rows.Close()

	rows.Next()

	err := rows.Scan(
		&id,
	)

	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (pg *PgDb) GetAuthorByName(b *v1.Book) (*v1.Author, error) {
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

	ok := rows.Next()

	if !ok {
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

func (pg *PgDb) CreateAuthorIfNew(b *v1.Book) (*int, error) {
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

func (pg *PgDb) BookCreate(b *v1.Book) (*int, error) {
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

	// Close the connection once we're done
	defer rows.Close()

	if ok := rows.Next(); !ok {
		return nil, nil
	}

	res := v1.GetEmptyBook()

	err = rows.Scan(
		&res.BookId,
		&res.CreatedTs,
		&res.Title,
		&res.PublishDate,
		&res.Edition,
		&res.Description,
		&res.Genre,
		&res.Author.AuthorId,
		&res.Author.CreatedTs,
		&res.Author.FirstName,
		&res.Author.LastName,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}
func (pg *PgDb) BookRemove(id int) error {
	return nil
}
func (pg *PgDb) BookFilter(b *v1.Book) ([]v1.Book, error) {
	return nil, nil
}
func (pg *PgDb) CollectionCreate(bookIds *[]int) (*int, error) {
	return nil, nil
}

func (pg *PgDb) CollectionGet(id int) (*v1.Collection, error) {
	return nil, nil
}
func (pg *PgDb) CollectionRemove(id int) error {
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
