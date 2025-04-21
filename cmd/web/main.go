package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/zachrundle/hotel-website/internal/config"
	"github.com/zachrundle/hotel-website/internal/handlers"
	"github.com/zachrundle/hotel-website/internal/models"
	"github.com/zachrundle/hotel-website/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {

	gob.Register(models.Reservation{})

	// change this to true when in production
	app.UseSSL = false
	app.SameSite = http.SameSiteLaxMode
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = app.SameSite
	session.Cookie.Secure = app.UseSSL

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalf("cannot create template cache: %v", err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
