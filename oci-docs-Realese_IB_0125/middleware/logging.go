package middleware

import (
	"bytes"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// LogRequest -
func LogRequest(log *zap.SugaredLogger, useLogging bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var blw *bodyLogWriter
		if useLogging {
			b, err := httputil.DumpRequest(c.Request, true)
			if err != nil {
				log.Errorf("DumpRequest error: %v", err)
				return
			}
			log.Debugf("Request Body: %s", string(b))
			blw = &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
		}

		c.Next()

		if useLogging {
			log.Debugf("Response Body: %s", blw.body.String())
		}
	}
}
