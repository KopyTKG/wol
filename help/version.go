package help

import "fmt"

var VERSION string = "v1.0.0"

func Version() {
    fmt.Printf("WoL %s\n", VERSION)
}
