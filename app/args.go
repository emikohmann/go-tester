package app

import (
    "flag"
    "errors"
    "fmt"
)

const (
    argConfig = "config"
)

type Arg struct {
    Key   string
    Value string
    Usage string
}

var (
    expectedArgs = []Arg{
        {
            Key:   argConfig,
            Value: "",
            Usage: "The config file",
        },
    }
    args = make(map[string]string)
)

func parseArgs() error {
    for i, arg := range expectedArgs {
        flag.StringVar(&expectedArgs[i].Value, arg.Key, arg.Value, arg.Usage)
    }
    flag.Parse()
    for _, arg := range expectedArgs {
        args[arg.Key] = arg.Value
    }
    for _, arg := range expectedArgs {
        if args[arg.Key] == "" {
            return errors.New(fmt.Sprintf("expected %s arg but not found", arg.Key))
        }
    }
    return nil
}
