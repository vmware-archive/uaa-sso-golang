package uaa

import (
    "crypto/tls"
    "io"
    "io/ioutil"
    "net/http"
)

// Http Client, wraps go's http.Client for our usecase
type Client struct {
    Host              string
    BasicAuthUsername string
    BasicAuthPassword string
    VerifySSL         bool
}

func NewClient(host, clientID, clientSecret string) Client {
    return Client{
        Host:              host,
        BasicAuthUsername: clientID,
        BasicAuthPassword: clientSecret,
    }
}

// Make request with the given basic auth and ssl settings, returns reponse code and body as a byte array
func (client Client) MakeRequest(method, path string, requestBody io.Reader) (int, []byte, error) {
    url := client.Host + path
    request, err := http.NewRequest(method, url, requestBody)
    if err != nil {
        return 0, nil, err
    }
    request.SetBasicAuth(client.BasicAuthUsername, client.BasicAuthPassword)
    request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    httpClient := http.Client{
        Transport: &http.Transport{
            TLSClientConfig: client.TLSConfig(),
        },
    }
    response, err := httpClient.Do(request)
    if err != nil {
        return 0, nil, err
    }

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return response.StatusCode, body, err
    }

    return response.StatusCode, body, nil
}

func (client Client) TLSConfig() *tls.Config {
    return &tls.Config{
        InsecureSkipVerify: !client.VerifySSL,
    }
}
