package main

import (
	"encoding/json"
	"fmt"

	"github.com/81120/tiny-parsec/ini"
)

func main() {
	// test json parser
	// str := `{ "key":  ["hello world", 123, true, null ] }, "val":  "kkkkk"}`
	// fmt.Println(str)
	// res := jsonp.JVal().Parse(str).Get()
	// s, _ := json.MarshalIndent(res.First, "", "  ")
	// fmt.Println(string(s))
	// fmt.Println(res.Second)

	// test ini parser
	str := `
	[section1]
	key1 = value1
	key2 = value2
	[section2]
	key3 = value3
	key4 = value4
	`
	fmt.Printf("%v", str)
	res := ini.ParseINI(str).Get()
	s, _ := json.MarshalIndent(res.First, "", "  ")
	fmt.Println(string(s), res.Second)

	// fmt.Println(strings.Split(str, "\n"))
}
