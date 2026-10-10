package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const maxLoggedBody = 2000

// responseRecorder captures the response body/status while still writing
// through to the real client — logging must never change what the caller
// receives.
type responseRecorder struct {
	http.ResponseWriter
	status       int
	bytesWritten int
	body         bytes.Buffer
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	n, err := r.ResponseWriter.Write(b)
	r.bytesWritten += n
	return n, err
}

// HTTPRequestLogger logs each request the same way chi's built-in
// middleware.Logger does — "METHOD scheme://host+uri PROTO" from clientIP -
// status bytesB in duration — then prints the request/response JSON bodies
// on the following lines, pretty-printed (2-space indent, like
// JSON.stringify(v, null, 2)) with sensitive fields (password, access_token,
// refresh_token...) redacted.
//
// Example:
//
//	2026/10/08 17:48:15 "POST http://localhost:8080/api/v1/auth/login HTTP/1.1" from 127.0.0.1:52292 - 200 295B in 881ms
//	  REQUEST  {
//	    "email": "a@b.com",
//	    "password": "***"
//	  }
//	  RESPONSE {
//	    "data": {
//	      "access_token": "***",
//	      ...
//	    }
//	  }
//
// REQUEST is always cyan; RESPONSE is colored by the response status
// (green/cyan/yellow/red), same as the status code on the summary line.
func HTTPRequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var reqBody []byte
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		rec := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)

		status := rec.status
		if status == 0 {
			status = http.StatusOK
		}

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		fmt.Printf("%s \"%s %s://%s%s %s\" from %s - %s%d%s %dB in %s\n",
			start.Format("2006/01/02 15:04:05"),
			r.Method, scheme, r.Host, r.RequestURI, r.Proto,
			r.RemoteAddr,
			statusColor(status), status, colorReset,
			rec.bytesWritten,
			time.Since(start),
		)

		if s := redactJSON(reqBody); s != "" {
			fmt.Printf("  %sREQUEST%s  %s\n", colorCyan, colorReset, truncate(s, maxLoggedBody))
		}
		if s := redactJSON(rec.body.Bytes()); s != "" {
			fmt.Printf("  %sRESPONSE%s %s\n", statusColor(status), colorReset, truncate(s, maxLoggedBody))
		}
	})
}

func statusColor(status int) string {
	switch {
	case status >= 500:
		return colorRed
	case status >= 400:
		return colorYellow
	case status >= 300:
		return colorCyan
	default:
		return colorGreen
	}
}
