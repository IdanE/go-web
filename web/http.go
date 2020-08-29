package web

import (
	"../template"
	"../util"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpMethod string
const(
	GET HttpMethod = "GET"
	POST HttpMethod = "POST"
	PUT HttpMethod = "PUT"
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
	Method HttpMethod
	Route string
	Handler func(w http.ResponseWriter, req *http.Request) Response
}

type Response struct {
	Text string
	Parameters []util.Pair
}

type HttpHandler struct {
	Routes []Route
}

func readFile(fileName string) (string, error) {
	fileContent, err := ioutil.ReadFile(fmt.Sprintf("./static/%s.html", fileName))
	if err != nil {
		// todo
		return "", err
	}
	return string(fileContent), nil
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	method, err := ValueOfHttpMethod(req.Method)
	uri := req.RequestURI
	if strings.HasPrefix(uri, "/resources") {
		uri = strings.ReplaceAll(uri, "/resources/", "")
		fileName := fmt.Sprintf("./static/resources/%s", uri)
		content, _ := ioutil.ReadFile(fileName)
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
				text, err = readFile(text)
				fmt.Fprintf(w, template.ProcessTemplate(text, response.Parameters...))
			} else {
				fmt.Fprintf(w, text)
			}

			return
		}
	}
}