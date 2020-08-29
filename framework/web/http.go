package web

import (
	"../template"
	"../util"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

var HttpMethods = []HttpMethod{
	GET,
	POST,
	PUT,
	DELETE,
}

func ValueOfHttpMethod(value string) (HttpMethod, error) {
	for _, httpMethod := range HttpMethods {
		if value == string(httpMethod) {
			return httpMethod, nil
		}
	}

	return "", errors.New(fmt.Sprintf("%s is not a valid HttpMethod", value))
}

type Route struct {
	Method  HttpMethod
	Route   string
	Handler func(w http.ResponseWriter, req *http.Request) Response
}

type Response struct {
	Text       string
	Parameters []util.Pair
}

type HttpHandler struct {
	Globals []util.Pair
	Routes  []Route
	Templates map[string]template.Template
}

func NewHttpHandler(templateDir string, routes []Route, globals []util.Pair) (*HttpHandler, error) {
	templates, err := template.LoadTemplates(templateDir)
	if err != nil {
		return nil, err
	}
	handler := HttpHandler{
		Globals: globals,
		Routes: routes,
		Templates: templates,
	}

	return &handler, nil
}

func (h HttpHandler) readTemplate(templateName string) (string, error) {
	if template, present := h.Templates[templateName]; present {
		return template.Content, nil
	}
	return "", errors.New(fmt.Sprintf("Could not find template %s", templateName))
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method, err := ValueOfHttpMethod(req.Method)
	uri := req.RequestURI
	if strings.HasPrefix(uri, "/static") {
		w.Header().Set("Content-Type", util.GetMimeType(filepath.Ext(uri)))
		content, _ := ioutil.ReadFile("./" + uri)
		w.Write(content)
		return
	}

	if len(uri) > 1 {
		uri = strings.TrimSuffix(uri, "/")

	}
	if err != nil {
		fmt.Fprint(w, "Invalid method")
		return
	}

	for _, route := range h.Routes {
		if uri == route.Route && method == route.Method {
			response := route.Handler(w, req)
			text := response.Text
			if strings.HasPrefix(text, "template:") {
				text = strings.ReplaceAll(text, "template:", "")
				text, err := h.readTemplate(text)
				if err != nil {
					fmt.Fprint(w, err)
					return
				}
				fmt.Fprintf(w, template.ProcessTemplate(text, append(response.Parameters, h.Globals...)...))
			} else {
				fmt.Fprintf(w, text)
			}

			return
		}
	}

	text, err := h.readTemplate("404")
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprintf(w, template.ProcessTemplate(text, append(
		h.Globals,
		util.Pair{Key: "title", Value: "Page Not Found"},
	)...))
}
