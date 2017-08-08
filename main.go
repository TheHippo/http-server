package main

import (
	"fmt"
	"net/http"
	"os"

	flags "github.com/jessevdk/go-flags"
)

type cliOptions struct {
	Port    uint   `short:"p" long:"port" description:"port to bind to" default:"8080"`
	Address string `short:"b" long:"address" description:"address to bind to" default:"localhost"`
	// Directory string `short:"d" long:"directory" description:"folder to serve" default:"."`
	Args struct {
		Directory string `positional-arg-name:"directory" description:"directory to serve"`
	} `positional-args:"yes"`
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
	fmt.Printf("Starting http-server at http://%s:%d\n", options.Address, options.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", options.Address, options.Port), fs)
}
