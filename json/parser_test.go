package json_test

import (
	"testing"

	"github.com/81120/tiny-parsec/json"
	"github.com/stretchr/testify/assert"
)

func TestParseString(t *testing.T) {
	t.Run("valid string", func(t *testing.T) {
		result := json.ParseJSON(`"hello"`)
		assert.Equal(t, "hello", result.Get().First.(json.JsonString).Val)
	})

	t.Run("invalid string", func(t *testing.T) {
		result := json.ParseJSON(`"unclosed`)
		assert.True(t, result.IsNothing())
	})
}

func TestParseArray(t *testing.T) {
	t.Run("nested arrays", func(t *testing.T) {
		result := json.ParseJSON(`[1, [true, null], "text"]`)
		assert.Len(t, result.Get().First.(json.JsonArray).Val, 3)
	})

	t.Run("empty array", func(t *testing.T) {
		result := json.ParseJSON(`[]`)
		assert.Len(t, result.Get().First.(json.JsonArray).Val, 0)
	})
}

func TestParseObject(t *testing.T) {
	t.Run("complex object", func(t *testing.T) {
		result := json.ParseJSON(`{
			"num": 42,
			"arr": [{"k": "v"}],
			"bool": false
		}`)
		m := result.Get().First.(json.JsonObject).Val
		assert.Equal(t, int64(42), m["num"].(json.JsonInt).Val)
		assert.Len(t, m["arr"].(json.JsonArray).Val, 1)
		assert.False(t, m["bool"].(json.JsonBool).Val)
	})

	t.Run("missing comma", func(t *testing.T) {
		result := json.ParseJSON(`{"a":1 "b":2}`)
		assert.True(t, result.IsNothing())
	})
}

func TestParseBoolean(t *testing.T) {
	t.Run("true value", func(t *testing.T) {
		result := json.ParseJSON(`true`)
		assert.True(t, result.Get().First.(json.JsonBool).Val)
	})

	t.Run("false value", func(t *testing.T) {
		result := json.ParseJSON(`false`)
		assert.False(t, result.Get().First.(json.JsonBool).Val)
	})
}
func TestParseInteger(t *testing.T) {
	t.Run("positive integer", func(t *testing.T) {
		result := json.ParseJSON(`42`)
		assert.Equal(t, int64(42), result.Get().First.(json.JsonInt).Val)
	})

	t.Run("negative integer", func(t *testing.T) {
		result := json.ParseJSON(`-42`)
		assert.Equal(t, int64(-42), result.Get().First.(json.JsonInt).Val)
	})
}
func TestParseFloat(t *testing.T) {
	t.Run("positive float", func(t *testing.T) {
		result := json.ParseJSON(`3.14`)
		assert.Equal(t, float64(3.14), result.Get().First.(json.JsonFloat).Val)
	})

	t.Run("negative float", func(t *testing.T) {
		result := json.ParseJSON(`-3.14`)
		assert.Equal(t, float64(-3.14), result.Get().First.(json.JsonFloat).Val)
	})
}
func TestParseNull(t *testing.T) {
	t.Run("null value", func(t *testing.T) {
		result := json.ParseJSON(`null`)
		assert.True(t, result.Get().First.(json.JsonNull).IsNil())
	})
}
