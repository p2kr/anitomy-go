package main

import (
	"encoding/json"
	"flag"
	"fmt"

	anitomy "github.com/p2kr/anitomy-go"
)

func main() {
	formatFlag := flag.String("format", "", "Output format (json)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		// Nothing to parse
		return
	}

	options := anitomy.DefaultOptions()

	if len(args) == 1 {
		elements := anitomy.ParseWithOptions(args[0], options)
		printResult(args[0], elements, *formatFlag)
	} else {
		results := anitomy.ParseBatch(args, options)
		for _, result := range results {
			printResult(result.Input, result.Elements, *formatFlag)
		}
	}
}

func printResult(input string, elements []anitomy.Element, format string) {
	if format == "json" {
		out := make(map[string]interface{})
		for _, e := range elements {
			kindStr := string(e.Kind)
			if existing, ok := out[kindStr]; ok {
				switch v := existing.(type) {
				case string:
					out[kindStr] = []string{v, e.Value}
				case []string:
					out[kindStr] = append(v, e.Value)
				}
			} else {
				out[kindStr] = e.Value
			}
		}
		bytes, _ := json.Marshal(out)
		fmt.Println(string(bytes))
	} else {
		fmt.Println(input)
		for _, e := range elements {
			fmt.Printf("%s: %s\n", e.Kind, e.Value)
		}
		fmt.Println()
	}
}
