package main

import (
    "github.com/emikohmann/go-tester/app"
)

const (
    defaultConfig = "./configs/config.json"
)

func main() {
    app.Start(defaultConfig)
}
