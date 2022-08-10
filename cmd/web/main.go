package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/joefazee/hnews/models"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	sessionKeyUserId   = "userId"
	sessionKeyUserName = "userName"
)

type application struct {
	appName string
	server  server
	debug   bool
	errLog  *log.Logger
	infoLog *log.Logger
	view    *jet.Set
	session *scs.SessionManager
	Models  models.Models
}

type server struct {
	host string
	port string
	url  string
}

func main() {

	migrate := flag.Bool("migrate", false, "should migrate - drop all tables")
	dsn := flag.String("dsn", "postgres://postgres:postgres@localhost/hnews?sslmode=disable", "postgres connection string")
	host := flag.String("host", "localhost", "domain name for the app")
	port := flag.String("port", "8009", "listening port")

	flag.Parse()

	server := server{
		host: *host,
		port: *port,
	}
	server.url = fmt.Sprintf("http://:%s:%s", *host, *port)

	db2, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db2.Close()

	// init upper/db
	upper, err := postgresql.New(db2)
	if err != nil {
		log.Fatal(err)
	}
	defer func(upper db.Session) {
		err := upper.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(upper)

	// run migration
	if *migrate {
		fmt.Println("Running migration")
		err = runMigrate(upper)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done running migration")
	}

	// init application
	app := &application{
		server:  server,
		appName: "HNews",
		debug:   true,
		infoLog: log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate|log.Lshortfile),
		errLog:  log.New(os.Stderr, "ERROR\t", log.Ltime|log.Ldate|log.Llongfile),
		Models:  models.New(upper),
	}

	// init jet template
	if app.debug {
		app.view = jet.NewSet(jet.NewOSFileSystemLoader("./views"), jet.InDevelopmentMode())
	} else {
		app.view = jet.NewSet(jet.NewOSFileSystemLoader("./views"))
	}

	// init session
	app.session = scs.New()
	app.session.Lifetime = 24 * time.Hour
	app.session.Cookie.Persist = true
	app.session.Cookie.Name = app.appName
	app.session.Cookie.Domain = app.server.host
	app.session.Cookie.SameSite = http.SameSiteStrictMode
	app.session.Store = postgresstore.New(db2)

	if err := app.listenAndServer(); err != nil {
		log.Fatal(err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrate(db db.Session) error {
	script, err := os.ReadFile("./migrations/tables.sql")
	if err != nil {
		return err
	}

	_, err = db.SQL().Exec(string(script))

	return err
}
