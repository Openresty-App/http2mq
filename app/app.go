package app

import (
	"log"
	"net/http"

	"fmt"

	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type App struct {
	S      *negroni.Negroni
	Router http.Handler
	Conf   *Http2MQ
}

func (a *App) Run() {
	//http.ListenAndServe(fmt.Sprintf(":%d", a.Conf.Port), a.Router)
	a.S.Run(fmt.Sprintf(":%d", a.Conf.Port))
}

func (a *App) Close() {
	a.Conf.AsyncProducer.Close()
	a.Conf.SyncProducer.Close()
}

func NewApp(f string) (*App, error) {
	c, err := InitConf(f)
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	r.Handle("/", NewKafka())

	n := negroni.New()
	n.Use(NewLogger())
	n.UseHandler(r)

	return &App{
		S:      n,
		Router: r,
		Conf:   c,
	}, nil
}

type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "[negroni] ", log.LstdFlags)}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	l.Printf("[%s %s] %v %s %v", r.Method, r.URL.Path, res.Status(),
		http.StatusText(res.Status()), time.Since(start))
}
