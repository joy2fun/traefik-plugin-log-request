package traefik_plugin_log_request

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

// Config holds the plugin configuration.
type Config struct {
	ResponseBody        bool   `json:"responseBody,omitempty"`
	RequestIDHeaderName string `json:"requestIDHeaderName,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type logRequest struct {
	name                string
	next                http.Handler
	responseBody        bool
	requestIDHeaderName string
}

type RequestData struct {
	URL          string `json:"url"`
	Host         string `json:"host"`
	Body         string `json:"body"`
	Headers      string `json:"headers"`
	ResponseBody string `json:"response_body"`
	RequestID    string `json:"request_id"`
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.RequestIDHeaderName == "" {
		config.RequestIDHeaderName = "X-Request-Id"
	}

	return &logRequest{
		name:                name,
		next:                next,
		responseBody:        config.ResponseBody,
		requestIDHeaderName: config.RequestIDHeaderName,
	}, nil
}

func (p *logRequest) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	requestId, _ := generateRandomID(16)

	if req.Header.Get(p.requestIDHeaderName) != "" {
		requestId = req.Header.Get(p.requestIDHeaderName)
	}

	req.Header.Set(p.requestIDHeaderName, requestId)

	body, err := io.ReadAll(req.Body)
	if err != nil {
	}

	req.Body = io.NopCloser(bytes.NewBuffer(body))

	wrappedWriter := &responseWriter{
		ResponseWriter: rw,
	}

	p.next.ServeHTTP(wrappedWriter, req)

	bodyBytes := wrappedWriter.buffer.Bytes()
	rw.Write(bodyBytes)

	headers := make(map[string]string)
	for name, values := range req.Header {
		headers[name] = values[0] // Take the first value of the header
	}

	jsonHeader, err := json.Marshal(headers)
	if err != nil {
	}

	requestData := RequestData{
		URL:       req.URL.String(),
		Host:      req.Host,
		Body:      string(body),
		Headers:   string(jsonHeader),
		RequestID: requestId,
	}

	if p.responseBody {
		responseBody := io.NopCloser(bytes.NewBuffer(bodyBytes))
		responseBodyBytes, err := io.ReadAll(responseBody)
		if err != nil {
		}

		requestData.ResponseBody = string(responseBodyBytes)
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
	}

	os.Stdout.WriteString(string(jsonData) + "\n")
}

type responseWriter struct {
	buffer bytes.Buffer

	http.ResponseWriter
}

func (r *responseWriter) Write(p []byte) (int, error) {
	return r.buffer.Write(p)
}

func (r *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("%T is not a http.Hijacker", r.ResponseWriter)
	}

	return hijacker.Hijack()
}

func (r *responseWriter) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func generateRandomID(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
