package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/SinmFull/BMS_prj/internal/data"
	"github.com/SinmFull/BMS_prj/internal/jsonlog"
	_ "github.com/go-sql-driver/mysql"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml")
	// viper.AddConfigPath(".")
	// if err := viper.ReadInConfig(); err != nil {
	// 	fmt.Printf("Error reading config file, %s", err)
	// }
	// fmt.Println("viper config:", viper.GetString("database.dsn"))

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "MySQL DSN")
	flag.Parse()

	// logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		// logger.Fatal(err)
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	// logger.Printf("database connection pool established")
	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.port),
		Handler:     app.routes(),
		IdleTimeout: time.Minute,
	}

	logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})
	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*sql.DB, error) {

	db, err := sql.Open("mysql", cfg.db.dsn)
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
