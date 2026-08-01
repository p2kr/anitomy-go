package anitomy

import "sync"

// Parse tokenizes and parses the input string with default options.
func Parse(input string) []Element {
	return ParseWithOptions(input, DefaultOptions())
}

// ParseWithOptions tokenizes and parses the input string with the provided options.
func ParseWithOptions(input string, options Options) []Element {
	tokenizer := NewTokenizer(input)
	tokenizer.Tokenize(options)

	parser := NewParser(tokenizer.Tokens)
	parser.Parse(options)

	return parser.Elements
}

// Result holds the parsed elements for a single input string.
type Result struct {
	Input    string
	Elements []Element
}

// ParseBatch parses multiple input strings concurrently.
func ParseBatch(inputs []string, options Options) []Result {
	results := make([]Result, len(inputs))
	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for i, input := range inputs {
		go func(index int, in string) {
			defer wg.Done()
			results[index] = Result{
				Input:    in,
				Elements: ParseWithOptions(in, options),
			}
		}(i, input)
	}

	wg.Wait()
	return results
}
