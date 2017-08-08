package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	flags "github.com/jessevdk/go-flags"
)

type cliOptions struct {
	Port    uint   `short:"p" long:"port" description:"port to bind to" default:"8080"`
	Address string `short:"b" long:"address" description:"address to bind to" default:"localhost"`
	Verbose bool   `short:"v" long:"verbose" description:"log every request"`
	// Directory string `short:"d" long:"directory" description:"folder to serve" default:"."`
	Args struct {
		Directory string `positional-arg-name:"directory" description:"directory to serve"`
	} `positional-args:"yes"`
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (st *statusWriter) WriteHeader(status int) {
	st.status = status
	st.ResponseWriter.WriteHeader(status)
}

func main() {
	var options cliOptions
	flags.Parse(&options)

	if options.Port == 0 && options.Address == "" {
		// -h was opened
		os.Exit(0)
	}

	if options.Args.Directory == "" {
		options.Args.Directory = "."
	}

	fs := http.FileServer(http.Dir(options.Args.Directory))
	if options.Verbose {
		fs = logHandler(fs.ServeHTTP)
	}
	abs, _ := filepath.Abs(options.Args.Directory)
	fmt.Printf("Serving %s at http://%s:%d\n", abs, options.Address, options.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", options.Address, options.Port), fs)
}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		sw := &statusWriter{
			ResponseWriter: w,
		}
		fn(sw, req)
		fmt.Printf("Served \"%s\" in %v [%d]\n", req.RequestURI, time.Now().Sub(start), sw.status)
	}
}
