package app

import (
    "fmt"
)

func Start(argConfig string) {
    config, err := LoadConfig(argConfig)
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
