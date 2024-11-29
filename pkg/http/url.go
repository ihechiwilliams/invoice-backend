package http

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	openSquareBracket    = "["
	closingSquareBracket = "]"
	emptySquareBrackets  = "[]"
)

var squareBracketsIndexRE = regexp.MustCompile(`\[\d+]`)

func StripIndexesFromArrayParams(req *http.Request) *http.Request {
	query := req.URL.Query()
	newQuery := make(url.Values)

	for key, values := range query {
		if strings.Contains(key, openSquareBracket) && strings.Contains(key, closingSquareBracket) {
			newKey := squareBracketsIndexRE.ReplaceAllString(key, emptySquareBrackets)

			newQuery[newKey] = values
		} else {
			newQuery[key] = values
		}
	}

	req.URL.RawQuery = newQuery.Encode()

	return req
}

func AddIndexesToArrayParams() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			newQuery := addIndexesToQueryWithArray(query)

			r.URL.RawQuery = newQuery.Encode()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func addIndexesToQueryWithArray(originalQuery url.Values) url.Values {
	newQuery := make(url.Values)

	for key, values := range originalQuery {
		if strings.Contains(key, emptySquareBrackets) {
			rawKey := strings.TrimSuffix(key, emptySquareBrackets)

			for index, value := range values {
				newQuery[rawKey+fmt.Sprintf("[%d]", index)] = []string{value}
			}
		} else {
			newQuery[key] = values
		}
	}

	return newQuery
}
