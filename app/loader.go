package app

import (
    "io/ioutil"
    "encoding/json"
)

type Payload map[string]interface{}

type Config struct {
    BaseURL             string    `json:"baseUrl"`
    Endpoints           []string  `json:"endpoints"`
    Methods             []string  `json:"methods"`
    Payloads            []Payload `json:"payloads"`
    RateLimiter         int       `json:"rate_limiter"`
    FilterResponseCodes []int     `json:"filter_response_codes"`
}

func LoadConfig(filename string) (*Config, error) {
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    var config Config
    if err := json.Unmarshal(bytes, &config); err != nil {
        return nil, err
    }
    return &config, nil
}
