package http

import (
	"bytes"
	"encoding/json"

	"golang.org/x/exp/slices"
)

const (
	ApplicationJSONType = "application/json"
	defaultFilterString = "filtered"
)

func compactJSON(src []byte) []byte {
	if !isValidJSON(src) {
		rawRespBody := &rawResponseBody{RawBody: string(src)}
		wrappedJSON, _ := json.Marshal(rawRespBody)

		return wrappedJSON
	}

	dst := &bytes.Buffer{}
	_ = json.Compact(dst, src)

	return dst.Bytes()
}

func isValidJSON(raw []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(raw, &js) == nil
}

type JSONRedactor struct {
	keysToHide       []string
	filterWithString string
}

type JSONRedactorOption func(*JSONRedactor)

func WithKeysToHide(keys []string) JSONRedactorOption {
	return func(jr *JSONRedactor) {
		jr.keysToHide = keys
	}
}

func WithFilterString(filter string) JSONRedactorOption {
	return func(jr *JSONRedactor) {
		jr.filterWithString = filter
	}
}

func NewJSONRedactor(options ...JSONRedactorOption) *JSONRedactor {
	jr := &JSONRedactor{}
	for _, opt := range options {
		opt(jr)
	}

	return jr
}

func (jr *JSONRedactor) HideJSONKeys(raw []byte) []byte {
	var payload interface{}

	unmarshalErr := json.Unmarshal(raw, &payload)
	if unmarshalErr != nil {
		return raw
	}

	redactedPayload := jr.redactKeys(payload)

	newData, marshalErr := json.Marshal(redactedPayload)
	if marshalErr != nil {
		return raw
	}

	return newData
}

func (jr *JSONRedactor) redactKeys(payload interface{}) interface{} {
	switch v := payload.(type) {
	case []interface{}:
		for i, element := range v {
			v[i] = jr.redactKeys(element)
		}
	case map[string]interface{}:
		for key, value := range v {
			if slices.Contains(jr.keysToHide, key) {
				v[key] = jr.filterWithString
			} else {
				v[key] = jr.redactKeys(value)
			}
		}
	}

	return payload
}
