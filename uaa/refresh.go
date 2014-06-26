package uaa

import (
    "encoding/json"
    "net/http"
    "net/url"
    "strings"
)

func Refresh(u UAA, refreshToken string) (Token, error) {
    token := NewToken()
    params := url.Values{
        "grant_type":    {"refresh_token"},
        "redirect_uri":  {u.RedirectURL},
        "refresh_token": {refreshToken},
    }
    code, body, err := u.makeRequest("POST", u.tokenURL(), strings.NewReader(params.Encode()))
    if err != nil {
        return token, err
    }
    switch {
    case code == http.StatusUnauthorized:
        return token, InvalidRefreshToken
    case code > 399:
        return token, NewFailure(code, body)
    }

    json.Unmarshal(body, &token)
    return token, nil
}
