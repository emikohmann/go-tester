package app

import (
    "fmt"
)

const (
    errLoadingConfig     = "Error loading config"
    errExecutingConfig   = "Error executing config"
    infExecutionSucceded = "Execution succeded"
)

func Start(argConfig string) {
    config, err := LoadConfig(argConfig)
    if err != nil {
        fmt.Println(errLoadingConfig, err)
        return
    }

    if err := config.Execute(); err != nil {
        fmt.Println(errExecutingConfig, err)
        return
    }

    fmt.Println(infExecutionSucceded)
}
