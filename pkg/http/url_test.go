package http

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripIndexesFromArrayParams(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://example.come?a[1]=bar&b[2]=foo&c[3]=qux&asd=123", http.NoBody)

	newReq := StripIndexesFromArrayParams(req)

	values := newReq.URL.Query()
	queryStr, _ := url.QueryUnescape(newReq.URL.RawQuery)

	assert.Equal(t, "a[]=bar&asd=123&b[]=foo&c[]=qux", queryStr)
	assert.Equal(t, []string{"bar"}, values["a[]"])
	assert.Equal(t, []string{"foo"}, values["b[]"])
	assert.Equal(t, []string{"qux"}, values["c[]"])
	assert.Equal(t, []string{"123"}, values["asd"])
}

func Test_addIndexesToQueryWithArray(t *testing.T) {
	originalQuery := url.Values{
		"a[0]": []string{"a 0 value"},
		"a[1]": []string{"a 1 value"},
		"b[]":  []string{"b 0 value", "b 1 value"},
		"c":    []string{"c value"},
	}

	newQuery := addIndexesToQueryWithArray(originalQuery)

	expectedQuery := url.Values{
		"a[0]": []string{"a 0 value"},
		"a[1]": []string{"a 1 value"},
		"b[0]": []string{"b 0 value"},
		"b[1]": []string{"b 1 value"},
		"c":    []string{"c value"},
	}

	assert.Equal(t, expectedQuery, newQuery)
}
