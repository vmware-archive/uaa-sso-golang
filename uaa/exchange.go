package uaa

import (
    "encoding/json"
    "net/url"
    "strings"
)

func Exchange(u UAA, authCode string) (Token, error) {
    token := NewToken()

    params := url.Values{
        "grant_type":   {"authorization_code"},
        "redirect_uri": {u.RedirectURL},
        "scope":        {u.Scope},
        "code":         {authCode},
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
