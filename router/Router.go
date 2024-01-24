package router

import (
	adminHandler "alwadimarkets/handler"
	"alwadimarkets/middleware"
	R "alwadimarkets/routes"
	"net/http"
	"strings"
)

type route struct {
	path    string
	handler middleware.RouteHandler
}

func Router(w http.ResponseWriter, r *http.Request) {
	var routes []route

	routes = append(routes, route{
		path:    "/",
		handler: middleware.WithHeaders(R.Home),
	})

	routes = append(routes, route{
		path:    "/product/:id",
		handler: middleware.WithHeaders(R.ProductId),
	})

	routes = append(routes, route{
		path:    "/category/:name",
		handler: middleware.WithHeaders(R.Category),
	})

	routes = append(routes, route{
		path:    "/company/:name",
		handler: middleware.WithHeaders(R.Company),
	})

	routes = append(routes, route{
		path:    "/subcategory/:name",
		handler: middleware.WithHeaders(R.SubCategory),
	})

	routes = append(routes, route{
		path:    "/offers",
		handler: middleware.WithHeaders(R.ProductInOffers),
	})

	routes = append(routes, route{
		path:    "/user/register",
		handler: middleware.WithHeaders(R.Register),
	})

	routes = append(routes, route{
		path:    "/user/login",
		handler: middleware.WithHeaders(R.Login),
	})

	routes = append(routes, route{
		path:    "/fav/add",
		handler: middleware.WithHeaders(R.AddFavourite),
	})

	routes = append(routes, route{
		path:    "/fav/delete",
		handler: middleware.WithHeaders(R.DeleteFavourite),
	})

	routes = append(routes, route{
		path:    "/fav",
		handler: middleware.WithHeaders(R.GetFavourite),
	})

	routes = append(routes, route{
		path:    "/profile",
		handler: middleware.WithHeaders(R.Profile),
	})

	routes = append(routes, route{
		path:    "/orderhistory",
		handler: middleware.WithHeaders(R.OrderHistory),
	})

	routes = append(routes, route{
		path:    "/order/:id",
		handler: middleware.WithHeaders(R.OrderPage),
	})

	routes = append(routes, route{
		path:    "/order",
		handler: middleware.WithHeaders(R.Order),
	})

	routes = append(routes, route{
		path:    "/delete",
		handler: middleware.WithHeaders(R.CancelOrder),
	})
	// Admin Routes

	routes = append(routes, route{
		path:    "/admin/login",
		handler: middleware.WithHeaders(adminHandler.AdminLogin),
	})
	routes = append(routes, route{
		path:    "/admin/products",
		handler: middleware.WithHeaders(adminHandler.GetProducts),
	})
	routes = append(routes, route{
		path:    "/admin/product/add",
		handler: middleware.WithHeaders(adminHandler.AddProduct),
	})
	routes = append(routes, route{
		path:    "/admin/product/edit",
		handler: middleware.WithHeaders(adminHandler.UpdateProduct),
	})
	routes = append(routes, route{
		path:    "/admin/product/delete",
		handler: middleware.WithHeaders(adminHandler.DeleteProduct),
	})
	routes = append(routes, route{
		path:    "/admin/orders",
		handler: middleware.WithHeaders(adminHandler.GetOrders),
	})
	routes = append(routes, route{
		path:    "/admin/order/edit",
		handler: middleware.WithHeaders(adminHandler.EditOrder),
	})
	routes = append(routes, route{
		path:    "/admin/order/:id",
		handler: middleware.WithHeaders(adminHandler.OrderPage),
	})
	routes = append(routes, route{
		path:    "/admin/users",
		handler: middleware.WithHeaders(adminHandler.GetUsers),
	})
	routes = append(routes, route{
		path:    "/admin/producttype",
		handler: middleware.WithHeaders(adminHandler.GetTypes),
	})
	routes = append(routes, route{
		path:    "/admin/producttype/add",
		handler: middleware.WithHeaders(adminHandler.AddTypes),
	})
	routes = append(routes, route{
		path:    "/admin/producttype/delete",
		handler: middleware.WithHeaders(adminHandler.DeleteTypes),
	})
	routes = append(routes, route{
		path:    "/admin/offers",
		handler: middleware.WithHeaders(adminHandler.GetOffers),
	})
	routes = append(routes, route{
		path:    "/admin/offers/add",
		handler: middleware.WithHeaders(adminHandler.AddOffers),
	})
	routes = append(routes, route{
		path:    "/admin/offers/delete",
		handler: middleware.WithHeaders(adminHandler.DeleteOffers),
	})
	routes = append(routes, route{
		path:    "/admin/subcategory",
		handler: middleware.WithHeaders(adminHandler.GetSubCategories),
	})
	routes = append(routes, route{
		path:    "/admin/subcategory",
		handler: middleware.WithHeaders(adminHandler.GetSubCategories),
	})
	routes = append(routes, route{
		path:    "/admin/subcategory/add",
		handler: middleware.WithHeaders(adminHandler.AddSubCategories),
	})
	routes = append(routes, route{
		path:    "/admin/subcategory/delete",
		handler: middleware.WithHeaders(adminHandler.DeleteSubCategories),
	})
	routes = append(routes, route{
		path:    "/admin/subcategory/edit",
		handler: middleware.WithHeaders(adminHandler.UpdateSubCategories),
	})
	routes = append(routes, route{
		path:    "/admin/companies",
		handler: middleware.WithHeaders(adminHandler.GetCompanies),
	})
	routes = append(routes, route{
		path:    "/admin/company/add",
		handler: middleware.WithHeaders(adminHandler.AddCompanies),
	})
	routes = append(routes, route{
		path:    "/admin/company/edit",
		handler: middleware.WithHeaders(adminHandler.EditCompanies),
	})
	routes = append(routes, route{
		path:    "/admin/company/delete",
		handler: middleware.WithHeaders(adminHandler.DeleteCompanies),
	})
	routes = append(routes, route{
		path:    "/admin/categories",
		handler: middleware.WithHeaders(adminHandler.GetCategories),
	})
	routes = append(routes, route{
		path:    "/admin/category/add",
		handler: middleware.WithHeaders(adminHandler.AddCategories),
	})
	routes = append(routes, route{
		path:    "/admin/category/edit",
		handler: middleware.WithHeaders(adminHandler.EditCategories),
	})
	routes = append(routes, route{
		path:    "/admin/category/delete",
		handler: middleware.WithHeaders(adminHandler.DeleteCategories),
	})

	for _, route := range routes {
		if matched, params := match(route.path, r.URL.Path); matched {
			route.handler(w, r, params)
			return
		}
	}

	http.Error(w, "Not Found", 404)
}

func match(pattern, path string) (bool, map[string]string) {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false, nil
	}

	if len(patternParts) == len(pathParts) {
		if pathParts[0] != patternParts[0] {
			return false, nil
		}
	}

	params := make(map[string]string)

	for i, patternPart := range patternParts {
		if strings.HasPrefix(patternPart, ":") {
			paramName := patternPart[1:]
			params[paramName] = pathParts[i]
		} else if patternPart != pathParts[i] {
			return false, nil
		}
	}

	return true, params
}
