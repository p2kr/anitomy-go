# Anitomy-Go

Anitomy-Go is a Golang port of [Anitomy](https://github.com/erengy/anitomy), a C++ library for parsing anime video filenames. 

## Features
- Ported directly from the original robust C++ parser.
- Extracts metadata such as title, episode, season, release group, resolution, and more.
- Concurrently process large batches of filenames with `anitomy.ParseBatch`.

## Installation

```sh
go get github.com/p2kr/anitomy-go
```

## Usage

Anitomy-Go provides several functions to fit different parsing needs: `Parse`, `ParseWithOptions`, and `ParseBatch`.

### 1. Basic Parsing

Use `anitomy.Parse()` to parse a single filename with the default options (which extracts all possible metadata).

```go
package main

import (
	"fmt"
	"github.com/p2kr/anitomy-go"
)

func main() {
	filename := "[SubsPlease] One Piece - 1080 (720p) [05B85B5E].mkv"
	
	// Parse with default options
	elements := anitomy.Parse(filename)
	
	for _, e := range elements {
		fmt.Printf("%s: %s\n", e.Kind, e.Value)
	}
}
```

### 2. Custom Parsing Options

Use `anitomy.ParseWithOptions()` when you want to customize which elements are parsed. Setting a property to `false` instructs the parser to skip that step, which can improve performance if you don't need all the data.

```go
options := anitomy.Options{
	ParseEpisode:         true,
	ParseEpisodeTitle:    false, // Skip parsing episode titles
	ParseFileChecksum:    false, // Skip CRC32 checksums
	ParseFileExtension:   true,
	ParsePart:            true,
	ParseReleaseGroup:    true,
	ParseSeason:          false, // Skip seasons
	ParseTitle:           true,
	ParseVideoResolution: true,
	ParseYear:            true,
}

elements := anitomy.ParseWithOptions("[Group] Anime - 01 (1080p).mkv", options)
```
*Note: `anitomy.DefaultOptions()` returns an `Options` struct with all flags set to `true`.*

### 3. Concurrent Batch Parsing

If you have a large dataset of filenames to parse, use `anitomy.ParseBatch()`. This function distributes the parsing workload across multiple Goroutines for significantly faster execution.

```go
filenames := []string{
	"[Group] Anime 1 - 01.mkv",
	"[Group] Anime 2 - 02.mkv",
	"[Group] Anime 3 - 03.mkv",
}

options := anitomy.DefaultOptions()
results := anitomy.ParseBatch(filenames, options)

for _, result := range results {
	fmt.Printf("Filename: %s\n", result.Input)
	for _, e := range result.Elements {
		fmt.Printf("  %s: %s\n", e.Kind, e.Value)
	}
}
```
## Tribute
This project does not ship with a functional CLI. For the original C++ CLI tool, see the [original Anitomy project](https://github.com/erengy/anitomy/blob/develop/docs/cli.md).

## License
Anitomy-Go is licensed under the Mozilla Public License 2.0 (MPL-2.0). See the `LICENSE` file for more details.
