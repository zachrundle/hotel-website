package handlers

import (
	"net/http"
	"log"
	"github.com/alexedwards/scs/v2"
	"github.com/zachrundle/hotel-website/internal/config"
	"encoding/gob"
	"time"
	"os"
	"github.com/zachrundle/hotel-website/internal/models"
	"github.com/zachrundle/hotel-website/internal/render"
	"github.com/justinas/nosurf"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"fmt"
	"path/filepath"

)


var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"


func getRoutes() http.Handler {
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.UseSSL = false
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	app.SameSite = http.SameSiteLaxMode
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = app.SameSite
	session.Cookie.Secure = app.UseSSL

	app.Session = session

	tc, err := CreateTestTemplateCache() 
	if err != nil {
		log.Fatalf("cannot create template cache: %v", err)
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := NewRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

  mux := chi.NewRouter()

	// middleware allows you to process a request and perform some logic on it (like restricting access to only authenticated users)

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)

	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.UseSSL,
		SameSite: app.SameSite,
	})

	return csrfHandler
}


// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	// the below method is better, but wanted to show another method above
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}

			myCache[name] = ts
		}
	}

	return myCache, nil
}