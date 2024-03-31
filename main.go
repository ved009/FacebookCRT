package main

import (
    "encoding/json"
    "fmt"
    "github.com/BurntSushi/toml"
    "io/ioutil"
    "net/http"
    "os"
)

type Config struct {
    Facebook struct {
        AccessToken string `toml:"access_token"`
    } `toml:"facebook"`
}

type CertificatesResponse struct {
    Data []struct {
        Domains []string `json:"domains"`
    } `json:"data"`
    Paging struct {
        Cursors struct {
            Before string `json:"before"`
            After  string `json:"after"`
        } `json:"cursors"`
    } `json:"paging"`
}

func fetchCertificates(url string) (*CertificatesResponse, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %v", err)
    }

    var certs CertificatesResponse
    if err := json.Unmarshal(body, &certs); err != nil {
        return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
    }

    return &certs, nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run script.go <query-domain>")
        return
    }

    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Error fetching user home directory:", err)
        return
    }

    configPath := homeDir + "/.Facebook.toml"

    var conf Config
    if _, err := toml.DecodeFile(configPath, &conf); err != nil {
        fmt.Println("Error reading config file:", err)
        return
    }

    accessToken := conf.Facebook.AccessToken
    queryDomain := os.Args[1]
    url := fmt.Sprintf("https://graph.facebook.com/v19.0/certificates?access_token=%s&fields=domains&limit=5000&pretty=0&query=%s", accessToken, queryDomain)

    for {
        certs, err := fetchCertificates(url)
        if err != nil {
            fmt.Println(err)
            break
        }

        for _, data := range certs.Data {
            for _, domain := range data.Domains {
                fmt.Println(domain)
            }
        }

        if certs.Paging.Cursors.After == "" {
            break
        }

        // Prepare URL for the next page
        url = fmt.Sprintf("%s&after=%s", url, certs.Paging.Cursors.After)
    }
}
