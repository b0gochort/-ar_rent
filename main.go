package main

import (
	"errors"
	"fmt"
	"github.com/b0gochort/car-rent/internal/handler"
	db "github.com/b0gochort/car-rent/internal/repository/postgres"
	"github.com/b0gochort/car-rent/model"
	"github.com/b0gochort/car-rent/pkg/logging"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
)

var Config model.Config

func main() {
	log := logging.InitLog()

	Config = *getConf()

	db, err := db.New(Config)
	if err != nil {
		log.Error("could not connect to db", "err", err)
		return
	}

	migrateDB(db, log)

	server := &http.Server{
		Addr:    Config.HTTP.Port,
		Handler: initRoutes(log, db),
	}

	log.Info("starting server on port", "port", Config.HTTP.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server", "err", err)
	}

}

func initRoutes(log *slog.Logger, db *sqlx.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("received request for", "path", r.URL.Path, "method", r.Method)

		switch r.URL.Path {
		case "/get-status":
			handler.CheckAuto(db, log).ServeHTTP(w, r)
		case "/get-price":
			handler.GetPrice(db, log).ServeHTTP(w, r)
		case "/create-session":
			handler.CreateRentSesion(db, log).ServeHTTP(w, r)
		case "/get-report":
			handler.GetCarUtilizationReport(db, log).ServeHTTP(w, r)

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
	})

	return mux
}

func getConf() *model.Config {
	viper.SetConfigFile("config/config.json")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}

	conf := &model.Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

func migrateDB(db *sqlx.DB, log *slog.Logger) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Error("Couldn't get database instance for running migrations: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "foj", driver)
	if err != nil {
		log.Error("Couldn't create migrate instance: ", err)
		return
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Error("Couldn't run database migration: %s", err.Error())
		return
	}

	log.Info("Database migration was run successfully")
	return
}
