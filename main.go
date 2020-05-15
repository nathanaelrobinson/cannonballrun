package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "database/sql"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
)

var (
	err     error
	jobChan chan Webhook
)

func main() {
	a := &App{}
	dbURI := os.Getenv("DB_CONNECTION")
	a.Initialize("mysql", dbURI)
	defer a.DB.Close()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// SET ROUTES HERE
	r := mux.NewRouter()
	r.Use(CommonMiddleware)
	r.Use(loggingMiddleware)
	r.HandleFunc("/register", a.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", a.LoginHandler).Methods("POST")
	r.HandleFunc("/home", home)
	r.HandleFunc("/authorize-strava", a.StravaAuthorization).Methods("GET", "POST")
	r.HandleFunc("/strava-webhook", a.CreateStravaWebhook).Methods("GET")
	r.HandleFunc("/strava-webhook", a.StravaWebhookHandler).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(JwtVerify)
	api.HandleFunc("/teams/{id:[0-9]+}/join", a.JoinTeam).Methods("GET")
	api.HandleFunc("/teams/{id:[0-9]+}", a.TeamHandlerDetail).Methods("GET", "PUT", "POST", "DELETE")
	api.HandleFunc("/teams", a.TeamHandlerList).Methods("GET")
	api.HandleFunc("/users", a.UserHandlerList).Methods("GET")
	api.HandleFunc("/users/me", a.UserHandlerDetail).Methods("GET")

	// Define a subrouter to handle files at static for accessing static content
	static := r.PathPrefix("/assets").Subrouter()
	static.Handle("/{*}/{*}", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	r.HandleFunc("/", index)

	// Logging for web server
	f, _ := os.Create("/var/log/golang/golang-server.log")
	defer f.Close()
	logger := handlers.CombinedLoggingHandler(f, r)

	// Logging for dev
	// logger := handlers.CombinedLoggingHandler(os.Stdout, r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	srv := &http.Server{
		Addr: ":" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      logger, // Pass our instance of gorilla/mux in.
	}

	// Set up a goroutine to handle asyncronous tasks
	jobChan = make(chan Webhook, 100)
	go a.UpdateWorkouts(jobChan)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Server Running on %q\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
