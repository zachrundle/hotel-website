package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/zachrundle/hotel-website/internal/config"
	"github.com/zachrundle/hotel-website/internal/handlers"
	"github.com/zachrundle/hotel-website/internal/helpers"
	"github.com/zachrundle/hotel-website/internal/models"
	"github.com/zachrundle/hotel-website/internal/render"
  "github.com/zachrundle/hotel-website/internal/driver"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

  defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})
  gob.Register(models.User{})
  gob.Register(models.Room{})
  gob.Register(models.Restriction{})

	// change this to true when in production
	app.UseSSL = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	app.SameSite = http.SameSiteLaxMode
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = app.SameSite
	session.Cookie.Secure = app.UseSSL

	app.Session = session

  // connect to database
  log.Println("Connecting to database...")
  db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=cybershady password=")
  if err != nil {
    log.Fatal("Cannot connect to database! Dying...", err)
  }
  log.Println("Connected to database!")



	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalf("cannot create template cache: %v", err)
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	render.NewRenderer(&app)
	return db, nil
}
