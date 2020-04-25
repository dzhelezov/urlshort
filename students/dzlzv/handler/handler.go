package urlshort

import (
	"log"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

type URLMapping struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectURL, ok := pathsToUrls[r.URL.Path]
		if !ok {
			log.Printf("No mapping found for %s\n", r.URL.Path)
			fallback.ServeHTTP(w, r)
			return
		}
		log.Printf("Redirecting to %s\n", redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	})
}

func parseYAML(yml []byte) (map[string]string, error) {
	ret := make(map[string]string)
	var parsed []URLMapping

	err := yaml.Unmarshal(yml, &parsed)
	if err != nil {
		return nil, err
	}

	for _, v := range parsed {
		log.Printf("%s: %s\n", v.Path, v.URL)
		ret[v.Path] = v.URL
	}

	return ret, nil
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedMap, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(parsedMap, fallback), nil
}
