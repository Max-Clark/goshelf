package db

import (
	"database/sql"
	"fmt"

	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
	_ "github.com/lib/pq"
)

type PgDb struct {
	SqlDb         *sql.DB
	SchemaVersion string // Used for migrations
	Config        ConnectionConfig
}

func (pg PgDb) BookCreate(v1.Book) error {

	return nil
}
func (pg PgDb) BookGet(id uint32) (*v1.Book, error) {
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
		&res.Publish_date,
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
func (pg PgDb) BookRemove(uint32) error
func (pg PgDb) BookFilter(v1.Book) ([]v1.Book, error)
func (pg PgDb) CollectionCreate([]uint32) error
func (pg PgDb) CollectionGet(uint32) (*v1.Collection, error)
func (pg PgDb) CollectionRemove(v1.Collection) error

func (pg PgDb) Connect() error {
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
