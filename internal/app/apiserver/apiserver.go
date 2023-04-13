package apiserver

import (
	"database/sql"
	"warehouse/internal/app/store/sqlstore"
)

func Start(config *Config) error {
	db, err := newDB("host=localhost database='warehouse' user='postgres' password='Tashkent123UzbMusayevSD' sslmode='disable'")
	// db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)
	srv := newServer(*store)
	srv.Logger.Info("address: ", config.s_address)
	return srv.Router.Run(config.s_address)
}

func newDB(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
