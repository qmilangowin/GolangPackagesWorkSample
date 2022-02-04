package main

import (
	"net/http"
)

//secureHeaders prevents xss attacks
func secureHeaders(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1:mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	}
}

func (app *ServerApplication) router(pattern string, next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if pattern == "/testinfo" {
			switch r.Method {
			case http.MethodGet:
				app.getTestInfo(w, r)
			case http.MethodPost:
				app.setTestInfo(w, r)
			default:
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			}
		} else {
			switch r.Method {
			case http.MethodGet:
				app.run(w, r)
			case http.MethodPost:
				app.run(w, r)
			default:
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			}

		}

	}
}
