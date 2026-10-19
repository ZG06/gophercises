package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"02-url-shortener/urlshort"
)

func main() {
	// -file flag to provide a custom yaml file
	yamlFile := flag.String("f", "urls.yaml", "a yaml file containing paths and target URLs")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Read data from the custom or default yaml file
	yaml, err := readYAML(*yamlFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

// Open the provided yaml file and return it's content in the form of bytes
func readYAML(filename string) ([]byte, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return f, nil
}
