package main

//
//
// import (
// 	"context"
// 	"database/sql"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"strconv"
// 	"sync/atomic"
// 	"syscall"
// 	"time"
//
// 	_ "github.com/mattn/go-sqlite3"
// )
//
// type middleware func(http.Handler) http.Handler
// type middlewares []middleware
//
// func (mws middlewares) apply(hdlr http.Handler) http.Handler {
// 	if len(mws) == 0 {
// 		return hdlr
// 	}
// 	return mws[1:].apply(mws[0](hdlr))
// }
//
// func (c *controller) shutdown(ctx context.Context, server *http.Server) context.Context {
// 	ctx, done := context.WithCancel(ctx)
//
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	go func() {
// 		defer done()
//
// 		<-quit
// 		signal.Stop(quit)
// 		close(quit)
//
// 		atomic.StoreInt64(&c.healthy, 0)
// 		server.ErrorLog.Printf("Server is shutting down...\n")
//
// 		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 		defer cancel()
//
// 		server.SetKeepAlivesEnabled(false)
// 		if err := server.Shutdown(ctx); err != nil {
// 			server.ErrorLog.Fatalf("Could not gracefully shutdown the server: %s\n", err)
// 		}
// 	}()
//
// 	return ctx
// }
//
// type controller struct {
// 	logger        *log.Logger
// 	nextRequestID func() string
// 	healthy       int64
// }
//
// var (
// 	db            *sql.DB
// 	err           error
// 	authenticated = true
// )
//
// func main() {
// 	// Connect to database for CRUD work
// 	db, err = sql.Open("sqlite3", "cbrun.db")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()
// 	// test connection
// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Launch and run the server including logging
// 	listenAddr := ":8000"
// 	if len(os.Args) == 2 {
// 		listenAddr = os.Args[1]
// 	}
//
// 	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
// 	logger.Printf("Server is starting...")
//
// 	c := &controller{logger: logger, nextRequestID: func() string { return strconv.FormatInt(time.Now().UnixNano(), 36) }}
// 	router := http.NewServeMux()
//
// 	// Register Routes
// 	router.HandleFunc("/", c.index)
// 	router.HandleFunc("/healthz", c.healthz)
// 	router.HandleFunc("/api", c.api)
// 	router.HandleFunc("/api/create", c.createWorkoutHandler)
//
// 	server := &http.Server{
// 		Addr:         listenAddr,
// 		Handler:      (middlewares{c.tracing, c.logging}).apply(router),
// 		ErrorLog:     logger,
// 		ReadTimeout:  5 * time.Second,
// 		WriteTimeout: 10 * time.Second,
// 		IdleTimeout:  15 * time.Second,
// 	}
// 	ctx := c.shutdown(context.Background(), server)
//
// 	logger.Printf("Server is ready to handle requests at %q\n", listenAddr)
// 	atomic.StoreInt64(&c.healthy, time.Now().UnixNano())
//
// 	if err := server.ListenAndServe(); err != http.ErrServerClosed {
// 		logger.Fatalf("Could not listen on %q: %s\n", listenAddr, err)
// 	}
// 	<-ctx.Done()
// 	logger.Printf("Server stopped\n")
//
// }
//
// // main_test.go, ensure that there is an index controller and healthz controller
// var (
// 	_ http.Handler = http.HandlerFunc((&controller{}).index)
// 	_ http.Handler = http.HandlerFunc((&controller{}).healthz)
// 	_ middleware   = (&controller{}).logging
// 	_ middleware   = (&controller{}).tracing
// )
