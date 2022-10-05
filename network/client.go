package network

import (
	"net/http"
	"time"
)

var transport = http.DefaultTransport.(*http.Transport).Clone()

func init() {
	transport.MaxIdleConns = 100
	transport.MaxIdleConnsPerHost = 100
	transport.MaxConnsPerHost = 100
	transport.IdleConnTimeout = 30 * time.Second
	transport.ResponseHeaderTimeout = 30 * time.Second
	transport.ExpectContinueTimeout = 30 * time.Second
}

var Client = &http.Client{
	Timeout:   10 * time.Second,
	Transport: transport,
}
