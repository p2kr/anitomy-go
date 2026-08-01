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

## Tribute
This project does not ship with a functional CLI. For the original C++ CLI tool, see the [original Anitomy project](https://github.com/erengy/anitomy/blob/develop/docs/cli.md).

## License
Anitomy-Go is licensed under the Mozilla Public License 2.0 (MPL-2.0). See the `LICENSE` file for more details.
