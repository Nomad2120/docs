package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model/context/auth"
)

// HTTPClient - интерфейс для http.Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ResourceError Объект для сообщения об ошибке
type ResourceError struct {
	URL        string
	StatusCode int
	Message    string
	Body       interface{}
	Err        error `json:"-"`
}

func (re *ResourceError) Error() string {
	return fmt.Sprintf(
		"Resource error: status code: %v,  err: %v",
		re.StatusCode,
		re.Err,
	)
	// return fmt.Sprintf(
	// 	"Resource error: URL: %s, status code: %v,  err: %v, body: %v",
	// 	re.URL,
	// 	re.StatusCode,
	// 	re.Err,
	// 	re.Body,
	// )
}

// RequestJSON  method(GET, POST, PUT, DELETE) return struct
func RequestJSON(ctx context.Context, client HTTPClient, method, url string, data []byte, headers map[string]string, response interface{}) (status int, respBody []byte, err error) {
	if headers == nil {
		headers = map[string]string{"Content-Type": "application/json"}
	} else {
		headers["Content-Type"] = "application/json"
	}

	if token := auth.FromContext(ctx); token != "" {
		headers["Authorization"] = token
	}

	status, respBody, err = send(ctx, client, method, url, data, headers)
	if err != nil {
		return
	}
	if response != nil && len(respBody) != 0 {
		err = json.Unmarshal(respBody, response)
	}
	return
}

func Request(ctx context.Context, client HTTPClient, method, url string, data []byte, headers map[string]string, response interface{}) (status int, respBody []byte, err error) {
	if token := auth.FromContext(ctx); token != "" {
		headers["Authorization"] = token
	}

	status, respBody, err = send(ctx, client, method, url, data, headers)
	if err != nil {
		return
	}
	if response != nil && len(respBody) != 0 {
		err = json.Unmarshal(respBody, response)
	}
	return
}

func send(ctx context.Context, client HTTPClient, method, urlString string, data []byte, headers map[string]string) (status int, buf []byte, err error) {
	request, err := http.NewRequest(method, urlString, bytes.NewBuffer(data))
	if err != nil {
		return status, nil, &ResourceError{URL: urlString, Err: err}
	}

	request = request.WithContext(ctx)

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	if strings.ContainsAny(urlString, "?") {
		urlTemp, err := url.Parse(urlString)
		if err != nil {
			return status, nil, &ResourceError{URL: urlString, Err: err}
		}
		urlQuery := urlTemp.Query()
		urlTemp.RawQuery = urlQuery.Encode()
		urlString = urlTemp.String()
	}

	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return status, nil, &ResourceError{URL: urlString, Err: err}
	}

	buf, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return status, nil, &ResourceError{URL: urlString, Err: err, StatusCode: response.StatusCode}
	}

	status = response.StatusCode
	if response.StatusCode > 399 {
		return status, buf, &ResourceError{
			URL:        urlString,
			Err:        fmt.Errorf("incorrect status code"),
			StatusCode: response.StatusCode,
			Message:    "incorrect response.StatusCode",
			Body:       string(data),
		}
	}

	return
}
