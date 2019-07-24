package app

import (
    "fmt"
    "sync"
    "net/http"
)

type Exploit struct {
    URL      string
    Methods  []string
    Payloads []Payload
}

type Potential struct {
    RequestMethod   string
    RequestURL      string
    RequestPayload  Payload
    ResponseStatus  int
    ResponseHeaders http.Header
    ResponsePayload []byte
}

type ExploitPotentials []Potential

func (config *Config) Execute() error {
    out := make(chan ExploitPotentials)
    go func() {
        for {
            exploitPotentials := <-out
            fmt.Println(exploitPotentials) // Change
        }
    }()

    domain := config.BuildURLs()

    var group sync.WaitGroup
    group.Add(len(domain))

    limiter := make(chan bool, config.RateLimiter)

    for _, url := range domain {
        limiter <- true
        exploit := &Exploit{
            URL:      url,
            Methods:  config.Methods,
            Payloads: config.Payloads,
        }
        go exploit.AsyncExecute(&group, limiter, out)
    }

    group.Wait()

    return nil
}

func (exploit *Exploit) AsyncExecute(group *sync.WaitGroup, limiter chan bool, out chan ExploitPotentials) {
    defer group.Done()
    out <- exploit.Execute()
    <-limiter
}

func (exploit *Exploit) Execute() ExploitPotentials {
    potentials := make(ExploitPotentials, 0)
    for _, method := range exploit.Methods {
        for _, payload := range exploit.Payloads {
            request := &Request{
                Method:  method,
                URL:     exploit.URL,
                Payload: payload,
            }
            response, apiErr := request.Do()
            if apiErr != nil {
                continue
            }
            potentials = append(potentials,
                Potential{
                    RequestMethod:   request.Method,
                    RequestURL:      request.URL,
                    RequestPayload:  request.Payload,
                    ResponseStatus:  response.StatusCode,
                    ResponseHeaders: response.Headers,
                    ResponsePayload: response.Payload,
                },
            )
        }
    }
    return potentials
}
