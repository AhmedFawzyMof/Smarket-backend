package router

import (
	"Smarket/routes"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

const (
	get    string = "GET"
	post   string = "POST"
	delete string = "DELETE"
	put    string = "PUT"
)

func Router(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler
	var slug string
	var id int

	path := r.URL.Path

	switch {
	case Match(path, "/"):
		handler = Get(routes.Home)
	case Match(path, "/product/([0-9]+)", &id):
		handler = Get(routes.ProductId{Id: id}.GetById)
	case Match(path, "/company/([^/]+)", &slug):
		handler = Get(routes.CompanySlug{Slug: slug}.GetBySlug)
	case Match(path, "/category/([^/]+)", &slug):
		handler = Get(routes.CategorySlug{Slug: slug}.GetBySlug)
	case Match(path, "/user/register"):
		handler = Post(routes.Register)
	case Match(path, "/user/login"):
		handler = Post(routes.Login)
	case Match(path, "/fav"):
		handler = Post(routes.Fav)
	case Match(path, "/fav/([^/]+)", &slug):
		handler = Get(routes.Favourite{Token: slug}.GetFav)
	case Match(path, "/Delfav"):
		handler = Post(routes.DelFav)
	default:
		http.NotFound(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}

func Match(path, pattern string, args ...interface{}) bool {
	regex := mustCompileCached(pattern)
	matches := regex.FindStringSubmatch(path)

	if len(matches) <= 0 {
		return false
	}

	for i, match := range matches[1:] {
		switch path := args[i].(type) {
		case *string:
			*path = match
		case *int:
			n, err := strconv.Atoi(match)
			if err != nil {
				return false
			}
			*path = n

		default:
			panic("args must be *string or *int")
		}

	}

	return true

}

func Get(handler http.HandlerFunc) http.HandlerFunc {
	return allowMethod(handler, get)
}

func Post(handler http.HandlerFunc) http.HandlerFunc {
	return allowMethod(handler, post)
}

func allowMethod(handler http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			w.Header().Set("Allow", method)
			http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

var (
	regexen = make(map[string]*regexp.Regexp)
	relock  sync.Mutex
)

func mustCompileCached(pattern string) *regexp.Regexp {
	relock.Lock()
	defer relock.Unlock()

	regex := regexen[pattern]
	if regex == nil {
		regex = regexp.MustCompile("^" + pattern + "$")
		regexen[pattern] = regex
	}
	return regex
}
