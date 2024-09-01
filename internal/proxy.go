package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CheckProxy verifies the given proxy string in the format "ip:port:username:password".
// It returns true if the proxy works, and false otherwise. The username and password can be null.
func CheckProxy(proxy string, testUrl string) bool {
	// Split the proxy string into parts
	parts := strings.Split(proxy, ":")
	if len(parts) < 2 || len(parts) > 4 {
		return false
	}

	ip := parts[0]
	port := parts[1]
	var username, password string
	if len(parts) == 4 {
		username = parts[2]
		password = parts[3]
	}

	// Create a proxy URL
	proxyURL := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", ip, port),
	}

	// Add authentication if username and password are provided
	if username != "" && password != "" {
		proxyURL.User = url.UserPassword(username, password)
	}

	// Create an HTTP client using the proxy
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// Make a request to a test URL to verify the proxy
	resp, err := client.Get(testUrl)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Return true if the status code is OK (200)
	return resp.StatusCode == http.StatusOK
}
