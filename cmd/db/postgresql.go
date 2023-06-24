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

func (pg PgDb) BookCreate(v1.Book) error
func (pg PgDb) BookGet(uint32) (*v1.Book, error)
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
