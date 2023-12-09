package server

import "net/http"

type handlerFunc func(http.ResponseWriter, *http.Request)

/*
TODO:
Refactor duplication later
*/
func Get(pattern string, handler handlerFunc) (string, handlerFunc) {
	return pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != pattern {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler(w, r)
	}
}

func Post(pattern string, handler handlerFunc) (string, handlerFunc) {
	return pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != pattern {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler(w, r)
	}
}
