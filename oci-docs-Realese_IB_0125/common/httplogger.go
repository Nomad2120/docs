package common

import (
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type loggedRoundTripper struct {
	rt  http.RoundTripper
	log HTTPLogger
}

func NewLoggedRoundTripper(rt http.RoundTripper,
	log HTTPLogger) *loggedRoundTripper {
	return &loggedRoundTripper{rt: rt, log: log}
}

func (lrt *loggedRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	lrt.log.logRequest(request)
	startTime := time.Now()
	response, err := lrt.rt.RoundTrip(request)
	duration := time.Since(startTime)
	lrt.log.logResponse(request, response, err, duration)
	return response, err
}

// HTTPLogger defines the interface to log http request and responses
type HTTPLogger interface {
	logRequest(*http.Request)
	logResponse(*http.Request, *http.Response, error, time.Duration)
}

type defaultLogger struct {
	dumpToFile bool
	log        *zap.SugaredLogger
}

// NewDefaultLogger -
func NewDefaultLogger(log *zap.SugaredLogger) *defaultLogger {
	save, _ := strconv.ParseBool(os.Getenv("dumpRequestsToFile"))
	return &defaultLogger{
		dumpToFile: save,
		log:        log,
	}
}

// LogRequest -
func (dl *defaultLogger) logRequest(req *http.Request) {
	buf, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		dl.log.Errorf("Error dump HTTP Request: %s", err.Error())
		return
	}
	dl.log.Debugf("HTTP Request: %s", string(buf))
	if dl.dumpToFile {
		saveToFile("HTTP Request:", buf)
	}
}

// LogResponse -
func (dl *defaultLogger) logResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	duration /= time.Millisecond
	if err != nil {
		dl.log.Debugf("HTTP Request method=%s host=%s path=%s duration=%d status=error error=%q", req.Method, req.Host, req.URL.Path, duration, err.Error())
	} else {
		buf, err := httputil.DumpResponse(res, true)
		if err != nil {
			dl.log.Debugf("Error dump HTTP Response: %s", err.Error())
			return
		}
		dl.log.Debugf("HTTP Response: %s duration: %d", string(buf), duration)
		if dl.dumpToFile {
			saveToFile("HTTP Response:", buf)
		}
	}
}

func saveToFile(title string, data []byte) error {
	fd, err := os.OpenFile("dump.log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()
	fd.WriteString(title + "\n" + string(data) + "\n\n\n\n")
	return nil
}
