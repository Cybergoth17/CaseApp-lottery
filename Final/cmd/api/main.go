package main

import (
	"Final/internal/data"
	"Final/internal/dbconnection"
	"Final/internal/jsonlogger"
	"Final/internal/mailer"
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"log"
	"os"
	"sync"
)

const version = "1.0.0"

type config struct {
	port    int
	env     string
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	logger *jsonlogger.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "f9773c0e408ce1", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "28b45d935c8415", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@almazsydykov768gmail.com>", "SMTP sender")
	flag.Parse()

	logger := jsonlogger.New(os.Stdout, jsonlogger.LevelInfo)

	res := dbconnection.DbConnection()
	db, er := sql.Open("postgres", res)
	if er != nil {
		log.Fatalf("postgres doesnt work : %s", er)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err := app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
