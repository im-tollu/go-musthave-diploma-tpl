package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/im-tollu/go-musthave-diploma-tpl/api"
	"github.com/im-tollu/go-musthave-diploma-tpl/config"
	"github.com/im-tollu/go-musthave-diploma-tpl/storage/auth/pg"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	conf, errConf := config.Load()
	if errConf != nil {
		log.Fatalf("Cannot load config: %s", errConf.Error())
	}

	log.Printf("Starting with config %#v", conf)

	if errMigrate := migrateDB(conf.DatabaseURI); errMigrate != nil {
		log.Fatalf("Cannot migrate DB: %s", errMigrate.Error())
	}

	db, errDB := newDataSource(conf.DatabaseURI)
	if errDB != nil {
		log.Fatalf("Cannot start DB: %s", errDB.Error())
	}

	authStorage, errAuthStorage := pg.NewAuthStorage(db)
	if errAuthStorage != nil {
		log.Fatalf("Cannot instantiate auth storage: %s", errAuthStorage.Error())
	}

	server, errServer := api.NewServer(conf.RunAddress, authStorage)
	if errServer != nil {
		log.Fatalf("Cannot start HTTP server: %s", errServer.Error())
	}

	awaitTermination()

	if errShutdown := server.Shutdown(context.Background()); errShutdown != nil {
		log.Fatalf("Could not gracefully stop the server: %s", errShutdown.Error())
	}
}

func awaitTermination() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
}

func migrateDB(databaseURL string) error {
	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("cannot init DB migrations: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("cannot apply migrations: %w", err)
	}

	return nil
}

func newDataSource(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot verify that DB connection is alive: %w", err)
	}

	return db, nil
}
