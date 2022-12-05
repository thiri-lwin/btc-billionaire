package postgres

import (
	"fmt"
	"log"
	"time"

	"btc_billionaire/utilities"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const dateFmt = "2006-01-02"

type Storage struct {
	db *sqlx.DB
}

type BTCInfo struct {
	ID      int             `db:"id"`
	Amount  decimal.Decimal `db:"amount"`
	Created time.Time       `db:"created"`
	Offset  int             `db:"offsettz"`
}

type BTCInfoResult struct {
	Amount  float64   `db:"amount" json:"amount"`
	Created time.Time `db:"created" json:"created"`
}

type Store interface {
	InsertBTCInfo(btcInfo BTCInfo) (BTCInfo, error)
	GetBTCHistory(startDate time.Time, endDate time.Time) ([]BTCInfoResult, error)
}

// create db if not exist
func createDB(config *utilities.Config) {
	dbstringWithoutDB := fmt.Sprintf(
		"user=%s host=%s port=%s sslmode=disable connect_timeout=5 TimeZone=UTC",
		config.DBUserName,
		config.DBHost,
		config.DBPort,
	)
	// If a password has been set, use it
	if password := config.DBPassword; password != "" {
		dbstringWithoutDB += " password=" + password
	}

	dbConn, err := newStorage(dbstringWithoutDB)
	if err != nil {
		log.Fatal(err)
	}

	dbConn.db.Exec("CREATE DATABASE " + config.DBName)

	dbConn.db.Close()
}

func NewStorageFromConfig(config *utilities.Config) (Store, error) {
	createDB(config)

	dbstring := fmt.Sprintf(
		"user=%s host=%s port=%s dbname=%s sslmode=disable connect_timeout=5 TimeZone=UTC",
		config.DBUserName,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	// If a password has been set, use it
	if password := config.DBPassword; password != "" {
		dbstring += " password=" + password
	}

	db, err := newStorage(dbstring)
	// run migration
	migrate(db)

	return db, err
}

func newStorage(dbString string) (*Storage, error) {
	db, err := sqlx.Connect("postgres", dbString)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to postgres '%s'", dbString)
	}
	return &Storage{db: db}, err
}

func (store *Storage) GetDBConn() *sqlx.DB {
	return store.db
}

func newTestStorageFromConfig(config *utilities.Config) (*Storage, func(), error) {
	dbstring := fmt.Sprintf(
		"user=%s host=%s port=%s dbname=%s sslmode=disable connect_timeout=5 TimeZone=UTC",
		config.DBUserName,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	// If a password has been set, use it
	if password := config.DBPassword; password != "" {
		dbstring += " password=" + password
	}
	db, teardown, err := newTestDB(dbstring, "../migrations")
	return &Storage{db: db}, teardown, err
}
