package template

import (
	"../util"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Template struct {
	Name string
	Content string
}


func LoadTemplates(path string) (map[string]Template, error) {
	templates := make(map[string]Template)
	matches, err := util.WalkMatch(path, "*.html")
	if err != nil {
		return nil, err
	}


	for _, match := range matches {
		content, err := ioutil.ReadFile(match)
		if err != nil {
			continue
		}
		match = strings.ReplaceAll(match, ".html", "")
		match = strings.Replace(match, fmt.Sprintf("%s/", filepath.Base(path)), "", 1)
		templates[match] = Template{
			Name: match,
			Content: string(content),
		}
	}

	return templates, nil
}