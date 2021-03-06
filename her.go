package her

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
)

const Version = "0.0.1 beta"

var (
	watcher *Watcher
	Config  *MergedConfig
)

type Application struct {
	Route    *Router
	Database *DB
	Template *TemplateFunc
}

func NewApplication(a ...interface{}) *Application {
	config := make(map[string]interface{})
	if len(a) > 0 {
		if v, ok := a[0].(map[string]interface{}); ok {
			config = v
		}
	}
	Config = loadConfig(config)
	application := &Application{
		Route:    newRouter(),
		Database: NewDB(),
		Template: newTemplateFunc(),
	}
	return application
}

func (app *Application) Start() {
	address := Config.String("Address")
	if address == "" {
		address = "0.0.0.0"
	}
	port := Config.Int("Port")
	if port == 0 {
		port = 8080
	}
	debug := Config.Bool("Debug")
	tmplPath := Config.String("TemplatePath")
	listen := fmt.Sprintf("%s:%d", address, port)

	templates = loadTemplate()

	watcher = NewWatcher()
	watcher.Listen(tmplPath)

	mux := http.NewServeMux()
	if debug {
		mux.Handle("/debug/pprof", http.HandlerFunc(pprof.Index))
		mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		mux.Handle("/debug/pprof/block", pprof.Handler("block"))
		mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	}
	mux.Handle("/", app.Route)

	l, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Listening on " + listen + "...")
	log.Fatal(http.Serve(l, mux))
}
