package urlshort

import (
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler that redirects urls if appropriate
func MapHandler(pathsToURL map[string]string, fallback http.Handler) http.HandlerFunc {
	return handler(pathsToURL, fallback)
}

func handler(pathsToURL map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL, ok := pathsToURL[r.URL.Path]
		if ok {
			redirect(w, r, redirectURL, 302)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url, code)
}

// YAMLHandler that parses input YAML and returns url.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathsToURL map[string]string = make(map[string]string)

	var yamlMapping []map[string]string
	err := yaml.Unmarshal(yml, &yamlMapping)
	if err != nil {
		fmt.Println(err)
		return func(w http.ResponseWriter, r *http.Request) {}, errors.New("Can't parse YAML")
	}
	length := len(yamlMapping)
	for i := 0; i < length; i++ {
		pathsToURL[yamlMapping[i]["path"]] = yamlMapping[i]["url"]
	}

	return handler(pathsToURL, fallback), nil
}
