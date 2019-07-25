package app

import (
    "fmt"
    "sync"
    "time"
    "net/http"
    "encoding/json"
    "github.com/emikohmann/go-tester/db"
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
    const (
        errSavingPotential = "error saving potential"
    )

    out := make(chan ExploitPotentials)

    go func() {
        for {
            exploitPotentials := <-out

            for _, potential := range exploitPotentials {
                if !potential.Match(config.FilterResponseCodes) {
                    continue
                }
                if err := potential.Save(); err != nil {
                    fmt.Println(errSavingPotential, err)
                    continue
                }
            }
        }
    }()

    fmt.Println("Building URLS..")
    domain := config.BuildURLs()

    var group sync.WaitGroup
    group.Add(len(domain))

    limiter := make(chan bool, config.RateLimiter)

    fmt.Println("Executing exploits...")

    for i, url := range domain {
        limiter <- true

        if i > 0 {
            fmt.Printf("\r%c[2K", 27)
        }

        exploit := &Exploit{
            URL:      url,
            Methods:  config.Methods,
            Payloads: config.Payloads,
        }

        fmt.Print(exploit.Methods, " >> ", exploit.URL)

        go exploit.AsyncExecute(&group, limiter, out)
    }

    fmt.Println()

    fmt.Println("Receiving results...")

    group.Wait()

    fmt.Println("Finishing execution...")

    time.Sleep(1 * time.Second)

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
        // Change payload validation
        for _, payload := range exploit.Payloads {
            request := &Request{
                Method:  method,
                URL:     exploit.URL,
                Payload: payload,
            }
            response, apiErr := request.Do()
            if apiErr != nil {
                // handle apiErr
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

func (potential *Potential) Match(responseCodes []int) bool {
    for _, responseCode := range responseCodes {
        if potential.ResponseStatus == responseCode {
            return false
        }
    }
    return true
}

func (potential *Potential) Save() error {
    const (
        potentialInsertQuery = "insert into potentials (request_method, request_url, request_payload, response_status, response_headers, response_payload) values (?, ?, ?, ?, ?, ?);"
    )
    requestPayload, err := json.Marshal(potential.RequestPayload)
    if err != nil {
        return err
    }
    responseHeaders, err := json.Marshal(potential.ResponseHeaders)
    if err != nil {
        return err
    }
    _, err = db.Client.Exec(
        potentialInsertQuery,
        potential.RequestMethod,
        potential.RequestURL,
        string(requestPayload),
        potential.ResponseStatus,
        string(responseHeaders),
        string(potential.ResponsePayload),
    )
    if err != nil {
        return err
    }
    return nil
}
