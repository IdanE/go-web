package main

import (
	"./framework/util"
	"./framework/web"
	"net/http"
)
func main() {
	globals := []util.Pair{
		{
			"siteName",
			"Site",
		},
		{
			"version",
			"1.0.0",
		},
	}
	routes := []web.Route{
		{
			Method: web.GET,
			Route:  "/",
			Handler: func(w http.ResponseWriter, r *http.Request) web.Response {
				return web.Response{
					Text: "template:index",
					Parameters: []util.Pair{{
						"name",
						"Test",
					}, {
						"title",
						"Home",
					}},
				}
			},
		},
		{
			Method: web.GET,
			Route:  "/about",
			Handler: func(w http.ResponseWriter, r *http.Request) web.Response {
				return web.Response{
					Text: "template:about",
					Parameters: []util.Pair{
						{
							"title",
							"About",
						},
					},
				}
			},
		},
	}
	handler, err := web.NewHttpHandler("./templates", routes, globals)
	if err != nil {
		panic(err)
	}
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}

