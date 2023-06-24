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

func (pg *PgDb) GetAuthorByName(b *v1.Book) (*v1.Author, error) {
	queryStr := fmt.Sprintf(`
		SELECT a.author_id, a.timestamp, a.first_name, a.last_name FROM %s.author a WHERE first_name = $1 AND last_name = $2
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

	rows.Next()

	a := v1.Author{}

	err = rows.Scan(
		&a.AuthorId,
		&a.Timestamp,
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
	`, pg.SchemaVersion)

	rows, err := pg.SqlDb.Query(
		queryStr,
		b.Author.FirstName,
		b.Author.LastName,
	)

	if err != nil {
		return err
	}

	defer rows.Close()

	return nil

}

func (pg *PgDb) BookCreate(b v1.Book) error {

	queryStr := fmt.Sprintf(`
		INSERT INTO %s.book (title,publish_date,edition,description,genre,author_id)
		VALUES              ($1,   $2,          $3,     $4,         $5,   $6       )
	`, pg.SchemaVersion)

	publishDate := b.PublishDate.Format(time.RFC3339)

	rows, err := pg.SqlDb.Query(
		queryStr,
		b.Title,
		publishDate,
		*b.Edition,
		*b.Description,
		*b.Genre,
		b.Author.AuthorId,
	)

	if err != nil {
		return err
	}

	rows.Close()

	return nil
}

func (pg *PgDb) BookGet(id uint32) (*v1.Book, error) {
	rows, err := pg.SqlDb.Query(`
			select b.book_id,b.title,b.publish_date,b.edition,b.description,b.genre,a.author_id,a.first_name,a.last_name
			from v1.book b 
			inner join v1.author a on b.author_id = a.author_id 
			where b.book_id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	// Close the connection once we're done
	defer rows.Close()

	res := v1.GetEmptyBook()

	rows.Next()

	err = rows.Scan(
		&res.BookId,
		&res.Title,
		&res.PublishDate,
		&res.Edition,
		&res.Description,
		&res.Genre,
		&res.Author.AuthorId,
		&res.Author.FirstName,
		&res.Author.LastName,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}
func (pg *PgDb) BookRemove(uint32) error {
	return nil
}
func (pg *PgDb) BookFilter(v1.Book) ([]v1.Book, error) {
	return nil, nil
}
func (pg *PgDb) CollectionCreate([]uint32) error {
	return nil
}

func (pg *PgDb) CollectionGet(uint32) (*v1.Collection, error) {
	return nil, nil
}
func (pg *PgDb) CollectionRemove(v1.Collection) error {
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
