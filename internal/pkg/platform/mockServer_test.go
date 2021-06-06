package platform

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

type countingServer struct {
	s          *httptest.Server
	successful int
	failed     []string
}

// Create a mock HTTP Server that will return a response with HTTP code and body.
func mockServer(code int, body string) *httptest.Server {
	serv := mockServerForQuery("", code, body)
	return serv.s
}

// mockServerForQuery returns a mock server that only responds to a particular query string.
func mockServerForQuery(query string, code int, body string) *countingServer {
	server := &countingServer{}


	server.s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(body,"{REPLACE_URL}"){
			body=strings.ReplaceAll(body,"{REPLACE_URL}",r.Host)
		}
		if query != "" && r.URL.RawQuery != query {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(query, r.URL.RawQuery, false)
			log.Printf("Query != Expected Query: %s", dmp.DiffPrettyText(diffs))
			server.failed = append(server.failed, r.URL.RawQuery)
			http.Error(w, "fail", 999)
			return
		}
		server.successful++

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		fmt.Fprintln(w, body)
	}))

	return server
}
