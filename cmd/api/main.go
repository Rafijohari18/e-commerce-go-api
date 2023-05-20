package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/models"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on ")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://postgres@localhost/gomovie_react?sslmode=disable", "Postgres connection config")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	fmt.Println("Server is running...")

	// http.HandleFunc("/status", func(rw http.ResponseWriter, r *http.Request) {
	// 	currentStatus := AppStatus{
	// 		Status:      "Online",
	// 		Environment: cfg.env,
	// 		Version:     version,
	// 	}

	// 	res, err := json.MarshalIndent(currentStatus, "", "\t")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	rw.Header().Set("Content-Type", "application/json")
	// 	rw.WriteHeader(http.StatusOK)
	// 	rw.Write(res)
	// })

	// err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), nil)

	serve := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting sever on port ", cfg.port)

	err = serve.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
