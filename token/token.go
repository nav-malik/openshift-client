package token

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func GetToken(username string, password []byte, serverURL string) (token string, err error) {

	if len(username) == 0 || username == "" || len(serverURL) == 0 || serverURL == "" || password == nil || len(password) < 1 {
		fmt.Println(username, password, serverURL)
		return "", errors.New("username, password, and serverURL are required")
	}
	// Create a new HTTP client that doesn't verify the server's TLS certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Set the client to not follow redirects
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	// Create a new GET request with the appropriate headers and authentication
	req, err := http.NewRequest("GET", serverURL+"/oauth/authorize?response_type=token&client_id=openshift-challenging-client", nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, string(password))
	req.Header.Set("X-CSRF-Token", "xxx")

	// Make the request and check for an error
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Get the Location header from the response
	location := resp.Header.Get("Location")

	// Extract the access token from the Location header
	locationURL, err := url.Parse(location)
	if err != nil {
		panic(err)
	}
	queryString := locationURL.Fragment
	decodedQueryString, err := url.QueryUnescape(queryString)
	if err != nil {
		panic(err)
	}
	keyValuePairs := strings.Split(decodedQueryString, "&")
	params := make(map[string]string)
	for _, pair := range keyValuePairs {
		splitPair := strings.Split(pair, "=")
		key := splitPair[0]
		value := splitPair[1]
		params[key] = value
	}
	accessToken := params["access_token"]

	return accessToken, nil
}
