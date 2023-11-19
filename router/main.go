package router

import (
	admin "Smarket/Admin"
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
	var table string

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
	case Match(path, "/subcategory/([^/]+)", &slug):
		handler = Get(routes.SubCategorySlug{Slug: slug}.GetBySlug)
	case Match(path, "/user/register"):
		handler = Post(routes.Register)
	case Match(path, "/user/login"):
		handler = Post(routes.Login)
	case Match(path, "/fav/add"):
		handler = Post(routes.Fav)
	case Match(path, "/fav"):
		handler = Post(routes.GetFav)
	case Match(path, "/Delfav"):
		handler = Post(routes.DelFav)
	case Match(path, "/profile/get"):
		handler = Post(routes.GetUserData)
	case Match(path, "/orderhistory"):
		handler = Post(routes.OrdersHistory)
	case Match(path, "/order"):
		handler = Post(routes.MakeOrders)
	case Match(path, "/order/delete"):
		handler = Post(routes.CancelOrder)
	case Match(path, "/offers"):
		handler = Get(routes.GetProductsOffers)
	case Match(path, "/foryou"):
		handler = Post(routes.ForYou)
	case Match(path, "/admin/products"):
		handler = Get(admin.GetProducts)
	case Match(path, "/admin/orders"):
		handler = Get(admin.GetOrders)
	case Match(path, "/admin/companies"):
		handler = Get(admin.GetComponies)
	case Match(path, "/admin/categories"):
		handler = Get(admin.GetCategories)
	case Match(path, "/admin/offers"):
		handler = Get(admin.GetOffers)
	case Match(path, "/admin/users"):
		handler = Get(admin.GetUsers)
	case Match(path, "/admin/subcategory"):
		handler = Get(admin.GetSubCategories)
	case Match(path, "/admin/([^/]+)/edit", &table):
		handler = Post(admin.Table{Name: table}.EditTable)
	case Match(path, "/admin/([^/]+)/delete", &table):
		handler = Post(admin.Table{Name: table}.DeleteTable)
	case Match(path, "/admin/([^/]+)/add", &table):
		handler = Post(admin.Table{Name: table}.CreateTable)
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
