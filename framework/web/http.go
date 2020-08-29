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
}

func readTemplate(fileName string) (string) {
	fileContent, err := ioutil.ReadFile(fmt.Sprintf("./templates/%s.html", fileName))
	if err != nil {
		// todo
		return ""
	}
	return string(fileContent)
}

func getMimeType(fileName string) string {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".css":
		return "text/css"
	default:
		return "text/plain"
	}
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method, err := ValueOfHttpMethod(req.Method)
	uri := req.RequestURI
	if strings.HasPrefix(uri, "/static") {
		w.Header().Set("Content-Type", getMimeType(uri))
		content, _ := ioutil.ReadFile("./" + uri)
		fmt.Fprintf(w, string(content))
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
				text = readTemplate(text)
				fmt.Fprintf(w, template.ProcessTemplate(text, append(response.Parameters, h.Globals...)...))
			} else {
				fmt.Fprintf(w, text)
			}

			return
		}
	}

	text := readTemplate("404")
	fmt.Fprintf(w, template.ProcessTemplate(text, append(
		h.Globals,
		util.Pair{Key: "title", Value: "Page Not Found"},
	)...))
}
