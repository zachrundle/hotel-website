package config

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog	  *log.Logger
	UseSSL        bool
	SameSite      http.SameSite
	Session *scs.SessionManager
}
