package main

import (
	"encoding/json"
	"fmt"

	"qgames/parser"
)

const filename = "qgames.log"

func main() {
	p := parser.Parser{}
	err := p.Parse(filename)
	if err != nil {
		panic(err)
	}

	out, _ := json.Marshal(p.Log)
	fmt.Println(string(out))
}
