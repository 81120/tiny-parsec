// Package json provides a set of parsers for JSON data using the tiny-parsec library.
package json

import (
	"github.com/81120/tiny-parsec/parser"
)

// JVal parses a JSON value, which can be a string, number, boolean, null, array, or object.
// It uses the OrElse combinator to try different parsers in order until one succeeds.
func JVal() parser.Parser[Json] {
	return parser.OrElse(
		JString(),
		JFloat(),
		JInt(),
		JBool(),
		JNull(),
		parser.Lazy(JArray),
		parser.Lazy(JObject),
	)
}

// JNull parses the JSON null value and returns a JsonNull object.
// It uses the Fmap combinator to transform the parsed string "null" into a JsonNull object.
func JNull() parser.Parser[Json] {
	return parser.Fmap(
		parser.Symbol("null"),
		func(_ string) Json {
			return JsonNull{}
		})
}

// JBool parses a JSON boolean value (true or false) and returns a JsonBool object.
// It uses the OrElse combinator to try parsing "true" or "false", and then the Fmap combinator to transform the result.
func JBool() parser.Parser[Json] {
	return parser.Fmap(
		parser.OrElse(parser.Symbol("true"), parser.Symbol("false")),
		func(str string) Json {
			return JsonBool{Val: str == "true"}
		})
}

// JInt parses a JSON integer value and returns a JsonInt object.
// It uses the Trim combinator to remove leading and trailing whitespace, and the Fmap combinator to transform the parsed integer.
func JInt() parser.Parser[Json] {
	return parser.Trim(
		parser.Fmap(parser.Integer(), func(i int64) Json {
			return JsonInt{Val: i}
		}))
}

// JFloat parses a JSON floating-point value and returns a JsonFloat object.
// It uses the Trim combinator to remove leading and trailing whitespace, and the Fmap combinator to transform the parsed float.
func JFloat() parser.Parser[Json] {
	return parser.Trim(
		parser.Fmap(parser.Float(), func(f float64) Json {
			return JsonFloat{Val: f}
		}))
}

// JString parses a JSON string value and returns a JsonString object.
// It uses the Trim combinator to remove leading and trailing whitespace, and the Fmap combinator to transform the parsed string.
func JString() parser.Parser[Json] {
	return parser.Trim(
		parser.Fmap(parser.String(), func(s string) Json {
			return JsonString{Val: s}
		}))
}

// JArray parses a JSON array value and returns a JsonArray object.
// It uses the Between combinator to parse the array enclosed in square brackets, and the SepBy combinator to parse the elements separated by commas.
func JArray() parser.Parser[Json] {
	return parser.Fmap(
		// 处理方括号包围的数组结构
		// Parse the array structure enclosed in square brackets
		parser.Between(
			parser.Trim(parser.Char('[')),                       // 左括号及空白
			parser.SepBy(JVal(), parser.Trim(parser.Char(','))), // 逗号分隔的元素
			parser.Trim(parser.Char(']')),                       // 右括号及空白
		),
		func(elements []Json) Json {
			return JsonArray{Val: elements}
		},
	)
}

// JPair parses a JSON key-value pair and returns a JsonPair object.
// It uses the Seq combinator to parse the key (a string), the colon separator, and the value, and then the Fmap combinator to transform the result.
func JPair() parser.Parser[JsonPair] {
	return parser.Fmap(
		parser.Seq(
			JString(),
			parser.Trim(
				parser.Fmap(
					parser.Char(':'),
					func(r rune) Json {
						return JsonString{Val: ":"}
					})),
			JVal()),
		func(tuple []Json) JsonPair {
			return JsonPair{
				Key:   tuple[0].(JsonString).Val,
				Value: tuple[2],
			}
		},
	)
}

// JObject parses a JSON object value and returns a JsonObject object.
// It uses the Between combinator to parse the object enclosed in curly braces, and the SepBy combinator to parse the key-value pairs separated by commas.
func JObject() parser.Parser[Json] {
	return parser.Fmap(
		parser.Between(
			parser.Trim(parser.Char('{')),
			parser.SepBy(JPair(), parser.Trim(parser.Char(','))),
			parser.Trim(parser.Char('}')),
		),
		func(pairs []JsonPair) Json {
			obj := make(map[string]Json)
			for _, pair := range pairs {
				obj[pair.Key] = pair.Value
			}
			return JsonObject{Val: obj}
		},
	)
}

func ParseJSON(jsonStr string) parser.ParserFuncRet[Json] {
	return JVal().Parse(jsonStr)
}
