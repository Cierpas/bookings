package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/cierpas/bookings/internal/config"
	"github.com/cierpas/bookings/internal/handlers"
	"github.com/cierpas/bookings/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager

const portNumber = ":8080"

func main() {

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	svr := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = svr.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
