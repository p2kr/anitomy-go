package anitomy_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	anitomy "github.com/p2kr/anitomy-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestData struct {
	Input  string                 `json:"input"`
	Output map[string]interface{} `json:"output"`
}

func TestAnitomyData(t *testing.T) {
	// Attempt to find data.json
	dataBytes, err := os.ReadFile(filepath.Join("..", "anitomy", "test", "data.json"))
	require.NoError(t, err, "failed to read data.json")

	var data []TestData
	err = json.Unmarshal(dataBytes, &data)
	require.NoError(t, err, "failed to parse data.json")

	for _, item := range data {
		t.Run(item.Input, func(t *testing.T) {
			elements := anitomy.Parse(item.Input)

			// Build a map of ElementKind string to array of values
			parsedMap := make(map[string][]string)
			for _, e := range elements {
				kindStr := string(e.Kind)
				parsedMap[kindStr] = append(parsedMap[kindStr], e.Value)
			}

			// Validate expected outputs against parsed outputs
			for expectedName, expectedVal := range item.Output {
				var expectedValues []string
				switch v := expectedVal.(type) {
				case string:
					expectedValues = append(expectedValues, v)
				case []interface{}:
					for _, item := range v {
						if s, ok := item.(string); ok {
							expectedValues = append(expectedValues, s)
						}
					}
				}

				parsedValues := parsedMap[expectedName]
				assert.Equal(t, expectedValues, parsedValues, "Mismatch for element: %s", expectedName)
			}

			// Validate there are no extra unexpected elements
			for parsedName, parsedValues := range parsedMap {
				_, expectedHas := item.Output[parsedName]
				if !expectedHas && len(parsedValues) > 0 {
					t.Errorf("Unexpected parsed element: %s = %v", parsedName, parsedValues)
				}
			}
		})
	}
}
