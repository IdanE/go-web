package template

import (
	"../util"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
)

func ProcessTemplate(template string, pairs ...util.Pair) string {
	template = replaceIncludes(template)
	template = replaceVariables(template, pairs...)
	return template
}

func replaceVariables(template string, pairs ...util.Pair) string {
	search := regexp.MustCompile("{{(\\s+)?([^\\s]+)(\\s+)?\\}\\}")

	result := search.ReplaceAllFunc([]byte(template), func(s []byte) []byte {
		value, err := getValueForVariable(search.ReplaceAllString(string(s), "$2"), pairs...)
		if err != nil {
			return []byte(fmt.Sprintf("{{ %s }}", err.Error()))
		}

		return []byte(value)
	})

	return string(result)
}

func replaceIncludes(template string) string {
	search := regexp.MustCompile("%include\\s(.+)")
	result := search.ReplaceAllFunc([]byte(template), func(s []byte) []byte {
		templateName := search.FindSubmatch(s)[1]
		content, err := ioutil.ReadFile(fmt.Sprintf("./templates/%s.html", templateName))
		if err != nil {
			return []byte("todo")
		}
		return []byte(search.ReplaceAllString(string(s), string(content)))
	})
	return string(result)
}

func getValueForVariable(variable string, pairs ...util.Pair) (string, error) {
	for _, pair := range pairs {
		if pair.Key == variable {
			return pair.Value, nil
		}
	}
	return "", errors.New(fmt.Sprintf("%s not found", variable))
}
