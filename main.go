package main

import (
	"encoding/json"
	"fmt"
)

const filename = "qgames.log"

func main() {
	p := Parser{}
	err := p.Parse(filename)
	if err != nil {
		panic(err)
	}

	out, _ := json.Marshal(p.Log)
	fmt.Println(string(out))
}
