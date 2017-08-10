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
	Cache   int    `short:"c" long:"cache" description:"set cache headers (-1 for no cache)" default:"0"`
	Args    struct {
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
	if options.Cache == -1 {
		fs = noCacheHandler(fs.ServeHTTP)
	}

	if options.Cache > 0 {
		fs = cacheHandler(options.Cache, fs.ServeHTTP)
	}

	abs, _ := filepath.Abs(options.Args.Directory)
	fmt.Printf("Serving %s at http://%s:%d\n", abs, options.Address, options.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", options.Address, options.Port), fs)
}

func cacheHandler(cacheDuration int, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", cacheDuration))
		fn(w, req)
	}
}

func noCacheHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Cache-Control", "no cache")
		req.Header.Del("If-Modified-Since")
		req.Header.Del("Cache-Control")
		fn(w, req)
	}
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
