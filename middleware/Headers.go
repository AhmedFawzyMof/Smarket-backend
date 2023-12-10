package middleware

import "net/http"

type RouteHandler func(http.ResponseWriter, *http.Request, map[string]string)

func WithHeaders(handler RouteHandler) RouteHandler {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		handler(w, r, params)
	}
}
