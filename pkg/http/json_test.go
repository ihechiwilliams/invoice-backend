package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JSONRedactor_HideJSONKeys(t *testing.T) {
	filteredKeys := []string{
		"secret1",
		"secret2",
		"secret3",
		"secret4",
		"secret5",
	}

	jr := NewJSONRedactor(
		WithKeysToHide(filteredKeys),
		WithFilterString(defaultFilterString),
	)

	t.Run("redacts keys from object", func(t *testing.T) {
		object := []byte(`{
	  "key1": "value",
	  "secret1": "value",
	  "key2": {
		"secret2": {
		  "secret3": {
			"key3": "value"
		  }
		}
	  },
	  "secret4": [
		{
		  "key4": "value",
		  "secret5": "abc"
		}
	  ]
	}`)
		filteredObject := jr.HideJSONKeys(object)

		assert.Equal(
			t,
			`{"key1":"value","key2":{"secret2":"filtered"},"secret1":"filtered","secret4":"filtered"}`,
			string(filteredObject),
		)
	})

	t.Run("redacts keys from an array", func(t *testing.T) {
		array := []byte(`[{
	  "key1": "value",
	  "secret1": "value",
	  "key2": {
		"secret2": {
		  "secret3": {
			"key3": "value"
		  }
		}
	  },
	  "secret4": [
		{
		  "key4": "value",
		  "secret5": "abc"
		}
	  ]
	}]`)
		filteredArray := jr.HideJSONKeys(array)

		assert.Equal(
			t,
			`[{"key1":"value","key2":{"secret2":"filtered"},"secret1":"filtered","secret4":"filtered"}]`,
			string(filteredArray),
		)
	})

	t.Run("skips invalid JSON", func(t *testing.T) {
		invalidJSON := []byte(`{"key1": "value"`)
		filteredJSON := jr.HideJSONKeys(invalidJSON)

		assert.Equal(t, `{"key1": "value"`, string(filteredJSON))
	})

	t.Run("skips simple string", func(t *testing.T) {
		simpleString := []byte(`"key1"`)
		filteredString := jr.HideJSONKeys(simpleString)

		assert.Equal(t, `"key1"`, string(filteredString))
	})
}
