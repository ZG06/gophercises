package urlshort

import (
	"encoding/json"
	"net/http"

	bolt "go.etcd.io/bbolt"
	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// Redirect to the destination URL (pathsToUrls value)
	// if the requested path was found in the pathsToUrls map.
	// Otherwise, delegates to the fallback handler.
	return func(w http.ResponseWriter, r *http.Request) {
		if dest, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, dest, 302)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// struct with path and url from the yaml data
type Data struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedYaml)

	return MapHandler(pathMap, fallback), nil
}

// The function unmarshals the yaml data in bytes
// and stores it into the Data slice.
func parseYAML(yml []byte) ([]Data, error) {
	d := []Data{}
	if err := yaml.Unmarshal(yml, &d); err != nil {
		return nil, err
	}

	return d, nil
}

func JSONHandler(js []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(js)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedJSON)

	return MapHandler(pathMap, fallback), nil
}

// The function unmarshals the json data in bytes
// and stores it into the Data slice.
func parseJSON(js []byte) ([]Data, error) {
	d := []Data{}
	if err := json.Unmarshal(js, &d); err != nil {
		return nil, err
	}

	return d, nil
}

func DBHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// path in the user's search bar
		pathToByte := []byte(r.URL.Path)
		var dest string = ""

		// check content of URLs bucket
		db.View(func(tx *bolt.Tx) error {
			// If the bucket is empty, return nothing.
			b := tx.Bucket([]byte("URLs"))
			if b == nil {
				return nil
			}

			// If the bucket isn't empty, get the url associated with the path
			// and save it to the dest string.
			val := b.Get(pathToByte)
			if val != nil {
				dest = string(val)
			}

			return nil
		})

		// If the destination url exists, redirect the user to it.
		if dest != "" {
			http.Redirect(w, r, dest, http.StatusFound)
			// Otherwise, delegate to the fallback handler.
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// The function build a map from the parsed data.
func buildMap(parsedData []Data) map[string]string {
	pathMap := make(map[string]string)

	for _, v := range parsedData {
		pathMap[v.Path] = v.Url
	}

	return pathMap
}
