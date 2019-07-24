package app

import (
    "fmt"
)

func (config *Config) BuildURLs() []string {
    const (
        urlFormat = "%s%s"
    )
    urls := make([]string, 0)
    for _, endpoint := range config.Endpoints {
        variants := Compose(endpoint)
        for _, variant := range variants {
            urls = append(urls, fmt.Sprintf(urlFormat, config.BaseURL, variant))
        }
    }
    return urls
}

func Compose(endpoint string) []string {
    return []string{
        endpoint,
    }
}
