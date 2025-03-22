package main

import (
	"fmt"

	"github.com/81120/tiny-parsec/json"
)

func main() {
	// test json parser
	str := `{ "key":  ["hello world", 123, true, null ] }, "val":  "kkkkk"}`
	res := json.JVal().Parse(str).Get()
	fmt.Println(res.First, res.Second)
}
