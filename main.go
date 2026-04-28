package main

import (
    "fmt"
    "os"

    "github.com/flessan/nurlaily/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m %v\n", err)
        os.Exit(1)
    }
}package main

import "github.com/flessan/nurlaily/cmd"

func main() {
    cmd.Execute()
}
