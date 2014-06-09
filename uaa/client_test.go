package uaa_test

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
    var client uaa.Client

    Describe("TLSConfig", func() {
        Context("when VerifySSL option is true", func() {
            It("uses a TLS config that verifies SSL", func() {
                client = uaa.NewClient("", "", "")
                client.VerifySSL = true
                Expect(client.TLSConfig().InsecureSkipVerify).To(BeFalse())
            })
        })

        Context("when VerifySSL option is false", func() {
            It("uses a TLS config that does not verify SSL", func() {
                client = uaa.NewClient("", "", "")
                client.VerifySSL = false
                Expect(client.TLSConfig().InsecureSkipVerify).To(BeTrue())
            })
        })
    })

    Describe("MakeRequest", func() {
        var server *httptest.Server
        var headers http.Header

        BeforeEach(func() {
            server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                w.WriteHeader(222)
                body := make([]byte, 40)
                n, err := req.Body.Read(body)
                if err != nil {
                    panic(err)
                }
                body = body[0:n]
                response := fmt.Sprintf("%s %s %s", req.Method, req.URL.Path, body)
                headers = req.Header
                w.Write([]byte(response))
            }))

        })

        It("makes an HTTP request with the given URL, HTTP method, and request body", func() {
            defer server.Close()

            client = uaa.NewClient(server.URL, "my-user", "my-pass")

            requestBody := strings.NewReader("key=value")
            code, body, err := client.MakeRequest("GET", "/something", requestBody)
            if err != nil {
                panic(err)
            }

            Expect(code).To(Equal(222))

            bodyText := string(body)
            Expect(bodyText).To(ContainSubstring("GET"))
            Expect(bodyText).To(ContainSubstring("/something"))
            Expect(bodyText).To(ContainSubstring("key=value"))
            Expect(headers["Content-Type"]).To(ContainElement("application/x-www-form-urlencoded"))
            Expect(strings.Join(headers["Authorization"], " ")).To(ContainSubstring("Basic"))
        })
    })
})
