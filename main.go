package main

import (
	"./template"
	"./util"
	"./web"
	"net/http"
)

func main() {
	template.ProcessTemplate("Hey")
	http.Handle("/", web.HttpHandler{
		Routes: []web.Route{
			{
				Method: web.GET,
				Route:  "/",
				Handler: func(w http.ResponseWriter, r *http.Request) web.Response {
					return web.Response{
						Text:       "template:index",
						Parameters: []util.Pair{{Key: "name", Value: "Test"}},
					}
				},
			},
			{
				Method: web.GET,
				Route:  "/about",
				Handler: func(w http.ResponseWriter, r *http.Request) web.Response {
					return web.Response{
						Text:       "template:about",
					}
				},
			},
		},
	})
	http.ListenAndServe(":8080", nil)
}
