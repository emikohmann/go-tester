package app

import (
    "fmt"
)

func Start() {
    if err := parseArgs(); err != nil {
        fmt.Println("Error parsing args", err)
        return
    }

    config, err := LoadConfig(args[argConfig])
    if err != nil {
        fmt.Println("Error loading config", err)
        return
    }

    if err := config.Execute(); err != nil {
        fmt.Println("Error executing config", err)
        return
    }

    fmt.Println("Execution succeded")
}
