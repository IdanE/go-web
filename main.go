package main

import (
	"./framework/util"
	"./framework/web"
	"net/http"
)

func main() {
	http.Handle("/", web.HttpHandler{
		Globals: []util.Pair{
			{
				"siteName",
				"Site",
			},
			{
				"version",
				"1.0.0",
			},
		},
		Routes: []web.Route{
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
		},
	})
	http.ListenAndServe(":8080", nil)
}
