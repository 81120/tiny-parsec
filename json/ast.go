// Package json defines a set of types to represent JSON data in Go.
package json

// Json is an interface that all JSON types must implement.
// It includes a method to indicate the type of the JSON value.
type Json interface {
	// jsonType is a method that all JSON types must implement.
	// It serves as a marker for the JSON type.
	jsonType()
}

// JsonNull represents a JSON null value.
type JsonNull struct{}

// IsNil returns true if the JsonNull value is nil.
func (j JsonNull) IsNil() bool {
	return true
}

// jsonType implements the Json interface for JsonNull.
func (j JsonNull) jsonType() {}

// JsonBool represents a JSON boolean value.
type JsonBool struct {
	// Val is the boolean value of the JSON boolean.
	Val bool
}

// jsonType implements the Json interface for JsonBool.
func (j JsonBool) jsonType() {}

// JsonInt represents a JSON integer value.
type JsonInt struct {
	// Val is the integer value of the JSON integer.
	Val int64
}

// jsonType implements the Json interface for JsonInt.
func (j JsonInt) jsonType() {}

// JsonFloat represents a JSON floating-point value.
type JsonFloat struct {
	// Val is the floating-point value of the JSON float.
	Val float64
}

// jsonType implements the Json interface for JsonFloat.
func (j JsonFloat) jsonType() {}

// JsonString represents a JSON string value.
type JsonString struct {
	// Val is the string value of the JSON string.
	Val string
}

// jsonType implements the Json interface for JsonString.
func (j JsonString) jsonType() {}

// JsonArray represents a JSON array value.
type JsonArray struct {
	// Val is the slice of Json values that make up the JSON array.
	Val []Json
}

// jsonType implements the Json interface for JsonArray.
func (j JsonArray) jsonType() {}

// JsonObject represents a JSON object value.
type JsonObject struct {
	// Val is the map of string keys to Json values that make up the JSON object.
	Val map[string]Json
}

// jsonType implements the Json interface for JsonObject.
func (j JsonObject) jsonType() {}

// JsonPair represents a key-value pair in a JSON object.
type JsonPair struct {
	// Key is the string key of the JSON pair.
	Key string
	// Value is the Json value of the JSON pair.
	Value Json
}
