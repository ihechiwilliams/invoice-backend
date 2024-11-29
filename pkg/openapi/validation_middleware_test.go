package openapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

func TestValidationMiddleware_Handler(t *testing.T) {
	loader := openapi3.NewLoader()
	doc, _ := loader.LoadFromFile("testdata/petstore.yml")

	vm := NewValidationMiddleware(
		WithDoc(doc),
		WithKinOpenAPIDefaults(),
	)

	vmHandler := vm.Handler()

	mux := http.DefaultServeMux

	mux.HandleFunc("/pets", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusCreated)
	})

	t.Run("when request body is empty", func(t *testing.T) {
		res := httptest.NewRecorder()
		buf := bytes.NewBuffer([]byte(`{ "data": {} }`))

		req := httptest.NewRequest(http.MethodPost, "https://example.com/pets", buf)

		req.Header.Set("Content-Type", "application/json")

		s := vmHandler(mux)

		s.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), `property \"age\" is missing`)
		assert.Contains(t, res.Body.String(), `property \"reference_id\" is missing`)
		assert.Contains(t, res.Body.String(), `property \"age\" is missing`)
	})

	t.Run("when uuid field has invalid format", func(t *testing.T) {
		res := httptest.NewRecorder()
		buf := bytes.NewBuffer([]byte(`{
		 "data": {
				"age": 5,
				"name": "Sparky",
				"reference_id": "non-valid-uuid"
			}
		}`))

		req := httptest.NewRequest(http.MethodPost, "https://example.com/pets", buf)
		req.Header.Set("Content-Type", "application/json")

		s := vmHandler(mux)

		s.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), `string doesn't match the format \"uuid\"`)
	})

	t.Run("when uuid field has default format", func(t *testing.T) {
		res := httptest.NewRecorder()
		buf := bytes.NewBuffer([]byte(`{
		 "data": {
				"age": 5,
				"name": "Sparky",
				"reference_id": "00000000-0000-0000-0000-000000000000"
			}
		}`))

		req := httptest.NewRequest(http.MethodPost, "https://example.com/pets", buf)
		req.Header.Set("Content-Type", "application/json")

		s := vmHandler(mux)

		s.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), `string doesn't match the format \"uuid\"`)
	})

	t.Run("when uuid field has valid format", func(t *testing.T) {
		res := httptest.NewRecorder()
		buf := bytes.NewBuffer([]byte(`{
		 "data": {
				"age": 5,
				"name": "Sparky",
				"reference_id": "7d339aad-ba59-47da-b762-8ffe6a37ed3b"
			}
		}`))

		req := httptest.NewRequest(http.MethodPost, "https://example.com/pets", buf)
		req.Header.Set("Content-Type", "application/json")

		s := vmHandler(mux)

		s.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
	})
}
