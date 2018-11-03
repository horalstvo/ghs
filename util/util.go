package util

import (
    "fmt"
    "os"
)

func Check(e error) {
    if e != nil {
        fmt.Printf("Error: %s\n", e.Error())
        os.Exit(1)
    }
}
