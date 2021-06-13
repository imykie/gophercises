package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"net/http"
)

type PathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/imyke":          "https://imyke.dev",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := YAMLHandler([]byte(yml), mapHandler)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	parsedMap := buildMap(parsedYaml)
	return MapHandler(parsedMap, fallback), nil
}

func parseYAML(yml []byte) ([]PathURL, error) {
	var pathUrls []PathURL
	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func buildMap(pathUrls []PathURL) map[string]string {
	m := make(map[string]string)
	for _, path := range pathUrls {
		m[path.Path] = path.URL
	}
	return m
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
