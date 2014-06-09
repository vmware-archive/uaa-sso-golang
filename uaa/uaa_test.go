package uaa_test

import (
    "net/http"
    "net/http/httptest"

    "github.com/pivotal-cf/go-uaa-client/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("UAA", func() {
    var auth uaa.UAA

    BeforeEach(func() {
        auth = uaa.NewUAA("http://login.example.com", "http://uaa.example.com", "the-client-id", "the-client-secret")
    })

    Describe("AuthorizeURL", func() {
        It("returns the URL for the /oauth/authorize endpoint", func() {
            Expect(auth.AuthorizeURL()).To(Equal("http://login.example.com/oauth/authorize"))
        })
    })

    Describe("LoginURL", func() {
        It("returns a url to be used as the redirect for authenticating with UAA", func() {
            auth.ClientID = "fake-client"
            auth.RedirectURL = "http://redirect.example.com"
            auth.Scope = "username,email"
            auth.State = "some-data"
            auth.AccessType = "offline"
            auth.ApprovalPrompt = "yes"
            expected := "http://login.example.com/oauth/authorize?access_type=offline&approval_prompt=yes&client_id=fake-client&redirect_uri=http%3A%2F%2Fredirect.example.com&response_type=code&scope=username%2Cemail&state=some-data"
            Expect(auth.LoginURL()).To(Equal(expected))
        })
    })

    Describe("Exchange", func() {
        var fakeUAAServer *httptest.Server

        Context("when UAA is responding normally", func() {
            BeforeEach(func() {
                fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                    if req.URL.Path == "/oauth/token" && req.Method == "POST" {
                        response := `{
                            "access_token": "access-token",
                            "refresh_token": "refresh-token",
                            "token_type": "bearer"
                        }`
                        w.WriteHeader(http.StatusOK)
                        w.Write([]byte(response))
                    } else {
                        w.WriteHeader(http.StatusNotFound)
                    }
                }))
                auth = uaa.NewUAA("http://login.example.com", fakeUAAServer.URL, "the-client-id", "the-client-secret")
            })

            AfterEach(func() {
                fakeUAAServer.Close()
            })

            It("returns a token received in exchange for a code from UAA", func() {
                token, err := auth.Exchange("1234")
                if err != nil {
                    panic(err)
                }

                Expect(token).To(Equal(uaa.Token{
                    Access:  "access-token",
                    Refresh: "refresh-token",
                }))
            })
        })

        Context("when UAA is not responding normally", func() {
            BeforeEach(func() {
                fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                    if req.URL.Path == "/oauth/token" && req.Method == "POST" {
                        w.WriteHeader(http.StatusUnauthorized)
                        w.Write([]byte(`{"errors": "Unauthorized"}`))
                    } else {
                        w.WriteHeader(http.StatusNotFound)
                    }
                }))
                auth = uaa.NewUAA("http://login.example.com", fakeUAAServer.URL, "the-client-id", "the-client-secret")
            })

            AfterEach(func() {
                fakeUAAServer.Close()
            })

            It("returns an error message", func() {
                _, err := auth.Exchange("1234")
                Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
                Expect(err.Error()).To(Equal(`UAA Failure: 401 {"errors": "Unauthorized"}`))
            })
        })
    })

    Describe("Refresh", func() {
        var fakeUAAServer *httptest.Server

        Context("when UAA is responding normally", func() {
            BeforeEach(func() {
                fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                    if req.URL.Path == "/oauth/token" && req.Method == "POST" {
                        err := req.ParseForm()
                        if err != nil {
                            panic(err)
                        }

                        if req.Form.Get("refresh_token") != "refresh-token" {
                            w.WriteHeader(http.StatusUnauthorized)
                            w.Write([]byte(`{"error":"invalid_token"}`))
                            return
                        }

                        response := `{
                            "access_token": "access-token",
                            "refresh_token": "refresh-token",
                            "token_type": "bearer"
                        }`

                        w.WriteHeader(http.StatusOK)
                        w.Write([]byte(response))
                    } else {
                        w.WriteHeader(http.StatusNotFound)
                    }
                }))
                auth = uaa.NewUAA("http://login.example.com", fakeUAAServer.URL, "the-client-id", "the-client-secret")
            })

            AfterEach(func() {
                fakeUAAServer.Close()
            })

            It("returns a token received in exchange for a refresh token", func() {
                token, err := auth.Refresh("refresh-token")
                Expect(err).To(BeNil())
                Expect(token.Access).To(Equal("access-token"))
            })

            It("returns an invalid refresh token error for invalid token", func() {
                _, err := auth.Refresh("bad-refresh-token")
                Expect(err).To(Equal(uaa.InvalidRefreshToken))
            })
        })

        Context("when UAA is not responding normally", func() {
            BeforeEach(func() {
                fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                    if req.URL.Path == "/oauth/token" && req.Method == "POST" {
                        response := `{"errors": "client_error"}`

                        w.WriteHeader(http.StatusMethodNotAllowed)
                        w.Write([]byte(response))
                    } else {
                        w.WriteHeader(http.StatusNotFound)
                    }
                }))
                auth = uaa.NewUAA("http://login.example.com", fakeUAAServer.URL, "the-client-id", "the-client-secret")
            })

            AfterEach(func() {
                fakeUAAServer.Close()
            })

            It("returns an error message", func() {
                _, err := auth.Refresh("refresh-token")
                Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
                Expect(err.Error()).To(Equal(`UAA Failure: 405 {"errors": "client_error"}`))
            })
        })
    })

    Describe("Refresh", func() {
        var fakeUAAServer *httptest.Server

        //"/oauth/token": {200, `{
        //"access_token": "some-client-auth-token",
        //"refresh_token": "refresh-token",
        //"token_type": "bearer"
        Context("when UAA is responding normally", func() {
            BeforeEach(func() {
                fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                    if req.URL.Path == "/oauth/token" && req.Method == "POST" {
                        err := req.ParseForm()
                        if err != nil {
                            panic(err)
                        }

                        if req.Form.Get("grant_type") != "client_credentials" {
                            w.WriteHeader(http.StatusNotAcceptable)
                            w.Write([]byte(`{"error":"unacceptable"}`))
                            return
                        }

                        response := `{
                            "access_token": "client-access-token",
                            "refresh_token": "refresh-token",
                            "token_type": "bearer"
                        }`

                        w.WriteHeader(http.StatusOK)
                        w.Write([]byte(response))
                    } else {
                        w.WriteHeader(http.StatusNotFound)
                    }
                }))
                auth = uaa.NewUAA("http://login.example.com", fakeUAAServer.URL, "the-client-id", "the-client-secret")
            })

            AfterEach(func() {
                fakeUAAServer.Close()
            })

            It("returns the client auth token", func() {
                token, err := auth.GetClientToken()
                Expect(err).To(BeNil())
                Expect(token.Access).To(Equal("client-access-token"))
            })
        })

        Context("when UAA is not responding normally", func() {
            BeforeEach(func() {
                fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                    if req.URL.Path == "/oauth/token" && req.Method == "POST" {
                        response := `{"errors": "Out to lunch"}`

                        w.WriteHeader(http.StatusGone)
                        w.Write([]byte(response))
                    } else {
                        w.WriteHeader(http.StatusNotFound)
                    }
                }))
                auth = uaa.NewUAA("http://login.example.com", fakeUAAServer.URL, "the-client-id", "the-client-secret")
            })

            AfterEach(func() {
                fakeUAAServer.Close()
            })

            It("returns an error message", func() {
                _, err := auth.GetClientToken()
                Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
                Expect(err.Error()).To(Equal(`UAA Failure: 410 {"errors": "Out to lunch"}`))
            })
        })
    })
})
