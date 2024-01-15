package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"qgames/parser"
)

func main() {
	// Define flags
	var inFile string
	var outFile string

	// Parse command-line arguments
	flag.StringVar(&inFile, "in", "qgames.log", "Input file name")
	flag.StringVar(&outFile, "out", "qgames.json", "Output file name")
	flag.Parse()

	p := parser.Parser{}
	parsedLog, err := p.Parse(inFile)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsedLog)
	if err := writeOutputToFile(outFile, parsedLog); err != nil {
		fmt.Println(parsedLog)
		panic(err)
	}

}

func writeOutputToFile(outFile, parsedLog string) error {
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.WriteString(f, parsedLog)
	if err != nil {
		return err
	}

	return nil
}
