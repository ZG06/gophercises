package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"02-url-shortener/urlshort"
)

func main() {
	// -f flag to provide a custom data file
	filePath := flag.String("f", "urls.yaml", "a data file containing paths and target URLs")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Read data from the custom or default data file
	data, err := readFile(*filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// check file extention of the provided data file
	ext := strings.ToLower(filepath.Ext(*filePath))
	var finalHandler http.HandlerFunc

	// choose the right handler for the extention
	switch ext {
	case ".json":
		finalHandler, err = urlshort.JSONHandler(data, mapHandler)
	case ".yaml", ".yml":
		finalHandler, err = urlshort.YAMLHandler(data, mapHandler)
	default:
		fmt.Printf("Unsupported file extension: %s. Please use .yaml, .yml, or .json\n", ext)
		os.Exit(1)
	}

	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", finalHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

// Open the provided data file and return its content in the form of bytes
func readFile(filename string) ([]byte, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return f, nil
}
