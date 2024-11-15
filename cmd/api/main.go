package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/abner-tech/Comments-Api.git/internal/data"
	"github.com/abner-tech/Comments-Api.git/internal/mailer"
	_ "github.com/lib/pq"
)

const appVersion = "3.0.0"

type serverConfig struct {
	port        int
	environment string
	db          struct {
		dsn string
	}
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

type applicationDependences struct {
	config       serverConfig
	logger       *slog.Logger
	commentModel data.CommentModel
	userModel    data.UserModel
	mailer       mailer.Mailer
	wg           sync.WaitGroup
	tokenModel   data.TokenModel
}

func main() {
	var settings serverConfig
	flag.IntVar(&settings.port, "port", 4000, "Server Port")
	flag.StringVar(&settings.environment, "env", "development", "Environment(development|staging|production)")
	//read the dsn
	flag.StringVar(&settings.db.dsn, "db-dsn", "postgres://comments:comments@localhost/comments?sslmode=disable", "PostgreSQL DSN")

	//limiter flags
	flag.Float64Var(&settings.limiter.rps, "limiter-rps", 2, "rate limiter maximum request per second")
	flag.IntVar(&settings.limiter.burst, "limiter-burst", 5, "rate limiter maximum burst")
	flag.BoolVar(&settings.limiter.enabled, "limiter-enabled", true, "enable rate limiter")

	//mailer flags
	flag.StringVar(&settings.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	//many ports are available, 25, 465, 587, 2525. If 25 doesn't work choose another
	flag.IntVar(&settings.smtp.port, "smtp-port", 25, "SMTP port")
	//personnal values provided by mailtrap
	flag.StringVar(&settings.smtp.username, "smtp-username", "839422506900bd", "SMTP username")
	flag.StringVar(&settings.smtp.password, "smtp-password", "ffb5cf13aa90aa", "SMTP password")
	flag.StringVar(&settings.smtp.sender, "smtp-sender", "Comments Community <no-reply@commentscommunity.amencias.net>", "SMTP sender")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	//the call to openDB() sets up our connection pool
	db, err := openDB(settings)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	//release the database connection before exiting
	defer db.Close()

	logger.Info("Database Connection Pool Established")

	appInstance := &applicationDependences{
		config:       settings,
		logger:       logger,
		commentModel: data.CommentModel{DB: db},
		userModel:    data.UserModel{DB: db},
		mailer:       mailer.New(settings.smtp.host, settings.smtp.port, settings.smtp.username, settings.smtp.password, settings.smtp.sender),
		tokenModel:   data.TokenModel{DB: db},
	}

	// apiServer := &http.Server{
	// 	Addr:         fmt.Sprintf(":%d", settings.port),
	// 	Handler:      appInstance.routes(),
	// 	IdleTimeout:  time.Minute,
	// 	ReadTimeout:  5 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// 	ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	// }

	// logger.Info("Starting Server", "address", apiServer.Addr, "environment", settings.environment, "limiter-enabled", settings.limiter)
	err = appInstance.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(settings serverConfig) (*sql.DB, error) {
	//open a connection pool
	db, err := sql.Open("postgres", settings.db.dsn)
	if err != nil {
		return nil, err
	}

	//set context to ensure DB operations dont take too long
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	//pinging connection pool to verify it was created, with a 5-second timeout
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	//return the connection pool (sql.DB)
	return db, nil
}
