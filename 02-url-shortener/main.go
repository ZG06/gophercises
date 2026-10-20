package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"02-url-shortener/urlshort"

	bolt "go.etcd.io/bbolt"
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

	// check file extention of the provided data file
	ext := strings.ToLower(filepath.Ext(*filePath))
	var finalHandler http.HandlerFunc

	// Read data from the custom or default data file

	// choose the right handler for the extention
	switch ext {
	case ".json":
		// Read data from a custom json file
		data, err := readFile(*filePath)
		if err != nil {
			fmt.Println("Error reading JSON file", err)
			os.Exit(1)
		}
		finalHandler, err = urlshort.JSONHandler(data, mapHandler)

	case ".yaml", ".yml":
		// Read data from a custom yaml file
		data, err := readFile(*filePath)
		if err != nil {
			fmt.Println("Error reading YAML file:", err)
			os.Exit(1)
		}
		finalHandler, err = urlshort.YAMLHandler(data, mapHandler)

	case ".db":
		// open a .db file
		db, err := bolt.Open(*filePath, 0600, nil)
		if err != nil {
			fmt.Println("Error reading BoltDB file", err)
			os.Exit(1)
		}

		defer db.Close()

		err = db.Update(func(tx *bolt.Tx) error {
			// open or create a URLs bucket
			b, err := tx.CreateBucketIfNotExists([]byte("URLs"))
			if err != nil {
				return err
			}

			// if the bucket is empty, add data from the pathsToUrls string
			if b.Stats().KeyN == 0 {
				for path, url := range pathsToUrls {
					err = b.Put([]byte(path), []byte(url))

					if err != nil {
						return err
					}
				}
			}

			return nil
		})

		finalHandler = urlshort.DBHandler(db, mapHandler)
	default:
		fmt.Printf("Unsupported file extension: %s. Please use .yaml, .yml, .json, or .db\n", ext)
		os.Exit(1)
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
