package uaa

import (
    "encoding/json"
    "net/url"
    "strings"
)

// Retrieves ClientToken from UAA server
func GetClientToken(u UAA) (Token, error) {
    token := NewToken()
    params := url.Values{
        "grant_type":   {"client_credentials"},
        "redirect_uri": {u.RedirectURL},
    }
    code, body, err := u.makeRequest("POST", u.tokenURL(), strings.NewReader(params.Encode()))
    if err != nil {
        return token, err
    }

    if code > 399 {
        return token, NewFailure(code, body)
    }

    json.Unmarshal(body, &token)
    return token, nil
}
