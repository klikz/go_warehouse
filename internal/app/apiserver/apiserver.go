package apiserver

import (
	"database/sql"
	"warehouse/internal/app/store/sqlstore"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)
	srv := newServer(*store)
	// if err := s.configureLogger(config); err != nil {
	// 	return err
	// }
	// s.configureRouter()

	// if err := s.ConfigureStore(); err != nil {
	// 	return err
	// }
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
