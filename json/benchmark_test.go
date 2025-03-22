package json_test

import (
	"strings"
	"testing"

	"github.com/81120/tiny-parsec/json"
)

func BenchmarkSimpleObject(b *testing.B) {
	data := `{"name":"John", "age":30, "active":true}`
	for i := 0; i < b.N; i++ {
		json.ParseJSON(data)
	}
}

func BenchmarkNestedStructure(b *testing.B) {
	data := `{"a":{"b":{"c":{"d":[1,2,{"e":3}]}}}}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.ParseJSON(data)
	}
}

func BenchmarkLargeArray(b *testing.B) {
	var sb strings.Builder
	sb.WriteString(`[`)
	for i := 0; i < 1000; i++ {
		if i > 0 {
			sb.WriteString(`,`)
		}
		sb.WriteString(`{"id":`)
		sb.WriteString(string(rune('0' + i%10)))
		sb.WriteString(`}`)
	}
	sb.WriteString(`]`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.ParseJSON(sb.String())
	}
}

func BenchmarkMixedTypes(b *testing.B) {
	data := `{
		"str": "value",
		"num": 3.14,
		"bool": true,
		"null": null,
		"arr": [1, "two", false],
		"obj": {"key": [{}]}
	}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.ParseJSON(data)
	}
}
