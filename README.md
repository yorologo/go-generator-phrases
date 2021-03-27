# Go-Generator-Phrases

This package generates random phrases based on the _Cards Against Humanity_ game, the game with **politically incorrect humor**.

## Supported Go versions & installation

Go-Generator-Phrases requires version 1.14 or higher of Go (Download Go)

Just use go get to download and install gearbox

    go get -u github.com/yorologo/go-generator-phrases

## Example

```go
package main

import (
    "fmt"
    "github.com/gogearbox/gearbox"
)

func main() {
    // Setup the generator
    g := generator.New()

    // Print in console a random phrase auto-generated
    fmt.Println(g.Generate())
}
```
