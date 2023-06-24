package postgresql

import (
	"database/sql"
	"fmt"
	"strings"

	v1 "github.com/Max-Clark/goshelf/cmd/model/v1"
	_ "github.com/lib/pq"
)

func (pg *PgDb) CollectionCreate(title *string, bookIds []int) (*string, error) {
	queryStr := fmt.Sprintf(`
		INSERT INTO %s.collection (title)
		VALUES ($1)
	`, pg.SchemaVersion)

	rows, err := pg.SqlDb.Query(queryStr, title)

	if err != nil {
		return nil, err
	}

	rows.Close()

	if len(bookIds) > 0 {
		queryStr = fmt.Sprintf(`
		INSERT INTO %s.collection_books (title,book_id)
		VALUES `, pg.SchemaVersion)

		values := make([]string, len(bookIds))
		for i := 0; i < len(bookIds); i++ {
			values[i] = " ( $1, $" + fmt.Sprint(i+2) + " ) " // $1 == title, so index + 2
		}

		queryStr += strings.Join(values, ",")

		// Create varargs for query function
		varArgs := make([]interface{}, len(bookIds)+1)

		varArgs[0] = title
		for i := 0; i < len(bookIds); i++ {
			varArgs[i+1] = bookIds[i]
		}

		rows, err := pg.SqlDb.Query(queryStr, varArgs...)

		if err != nil {
			return nil, err
		}

		rows.Close()
	}

	return title, nil
}

func (pg *PgDb) CollectionGet(title *string) (*v1.Collection, error) {
	if title == nil {
		return nil, nil
	}

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
		WHERE cb.title = $1
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

func (pg *PgDb) CollectionRemove(title *string) error {
	if title == nil {
		return nil
	}

	queryStr := fmt.Sprintf(`
		DELETE FROM %s.collection c
		WHERE c.title = $1
	`, pg.SchemaVersion)

	rows, err := pg.SqlDb.Query(queryStr, *title)

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
