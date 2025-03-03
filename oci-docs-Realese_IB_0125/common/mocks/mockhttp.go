package mocks

import (
	"fmt"
	"net/http"
	"sync"
)

// RoundTripFunc - roudtrip function type
type RoundTripFunc func(req *http.Request, resp *http.Response) (*http.Response, error)

//EqualRequestFunc - function type to check request equivalence
type EqualRequestFunc func(req *http.Request) bool

//ResponseHandler - response handler type
type ResponseHandler struct {
	Response     *http.Response
	Handler      RoundTripFunc
	URL          string
	EqualRequest EqualRequestFunc
}

// RoundTripper - structure for emulating a response from a service
type RoundTripper struct {
	handlers  []ResponseHandler
	callCount int
	lock      sync.Mutex
}

// NewRoundTripper - creating a mock to simulate responses from http requests
func NewRoundTripper(rh []ResponseHandler) *RoundTripper {
	return &RoundTripper{handlers: rh, callCount: 0}
}

func defaultHandler(req *http.Request, resp *http.Response) (*http.Response, error) {
	return resp, nil
}

//RoundTrip - implements the RoundTripper interface.
func (m *RoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	m.lock.Lock()
	defer func() {
		m.callCount++
		m.lock.Unlock()
	}()

	if len(m.handlers) == 0 {
		return nil, fmt.Errorf("response handlers is empty")
	}
	if m.callCount > len(m.handlers)-1 {
		return nil, fmt.Errorf("invalid index %d, handlers array size %d", m.callCount, len(m.handlers))
	}

	n := -1
	if m.handlers[0].EqualRequest != nil {
		for i, resp := range m.handlers {
			if resp.EqualRequest(request) {
				n = i
				break
			}
		}

	} else {
		n = m.callCount
	}

	if n == -1 {
		return nil, fmt.Errorf(request.URL.Path + " not found in mock")
	}

	if m.handlers[n].Handler == nil {
		response, err := defaultHandler(request, m.handlers[n].Response)
		return response, err
	}

	response, err := m.handlers[m.callCount].Handler(request, m.handlers[n].Response)
	return response, err
}
