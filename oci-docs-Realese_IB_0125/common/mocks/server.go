package mocks

import (
	"io"
	"net/http"
	"net/http/httptest"
)

// DoRequest -
func DoRequest(r http.Handler, method, path string, username, password string, headers map[string]string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
