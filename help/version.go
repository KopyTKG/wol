package help

import "fmt"

var VERSION string = "v0.0.1"

func Version() {
    fmt.Printf("WoL %s\n", VERSION)
}
