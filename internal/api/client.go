package api

import (
	"net/http"
	"time"
)

const (
	// DefaultTimeout is the timeout for HTTP requests
	DefaultTimeout = 10 * time.Second
)

// DefaultClient is the HTTP client used for API requests
var DefaultClient = &http.Client{
	Timeout: DefaultTimeout,
}
