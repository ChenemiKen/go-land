package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/chenemiken/goland/bookings/helpers"
	"github.com/chenemiken/goland/bookings/internal/config"
	"github.com/chenemiken/goland/bookings/internal/drivers"
	"github.com/chenemiken/goland/bookings/internal/handlers"
	"github.com/chenemiken/goland/bookings/internal/models"
	"github.com/chenemiken/goland/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	listenForMail()

	fmt.Printf((fmt.Sprintf("Starting application on port %s \n", portNumber)))
	// _ = http.ListenAndServe(portNumber, nil)

	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err)
}

func run() (*drivers.DB, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	inProduction := flag.Bool("production", true, "App is in production")
	useCache := flag.Bool("cache", true, "App should use template cache")
	dbHost := flag.String("dbhost", "localhost", "database host")
	dbPort := flag.String("dbport", "5432", "database port")
	dbName := flag.String("dbname", "", "database name")
	dbUser := flag.String("dbuser", "", "database user")
	dbPass := flag.String("dbpass", "", "database pass")
	dbSSL := flag.String("dbssl", "disable",
		"database ssl settings (disable, prefer, required)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing database credential flags")
		os.Exit(1)
	}

	app.InProduction = *inProduction
	app.UseCache = *useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = *scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.Session = &session

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s "+
		"password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass,
		*dbSSL)
	db, err := drivers.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Can not connect to db... Dying.")
		return nil, err
	}
	app.InfoLog.Println("Connected to the database!")

	render.NewRenderer(&app)

	repo := handlers.NewRepo(db, &app)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	return db, nil
}
