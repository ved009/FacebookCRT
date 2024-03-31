package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
)

// Define structures to unmarshal the JSON data.
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

    queryDomain := os.Args[1]
    accessToken := "Access Token Here" // Replace with your actual access token
    baseurl := "https://graph.facebook.com/v19.0/certificates?access_token=%s&fields=domains&limit=5000&pretty=0&query=%s"
    url := fmt.Sprintf(baseurl, accessToken, queryDomain)

    for {
        certs, err := fetchCertificates(url)
        if err != nil {
            fmt.Println(err)
            break
        }

        // Process and print each domain on a new line
        for _, data := range certs.Data {
            for _, domain := range data.Domains {
                fmt.Println(domain)
            }
        }

        if certs.Paging.Cursors.After == "" {
            //fmt.Println("No more data available.")
            break
        }

        url = fmt.Sprintf("%s&after=%s", url, certs.Paging.Cursors.After)
    }
}
