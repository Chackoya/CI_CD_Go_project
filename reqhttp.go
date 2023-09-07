/*
- reqhttp.go
Functions in this file take care of making HTTP GET requests and error handling.

- AUTH requests if a token is provided. Can be provided in a .env file, please checkout .env.example template, or when using the docker image (deployed) by using the cmd:
docker run -e GITHUB_TOKEN=private_token USERNAME/IMAGE_NAME:TAG



- UNAUTH requests (limited) if no token is provided.




*/
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Utility function to make Requests (Authorized ones if token is defined on .env file)
func makeGETRequest(url, token string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// If token is present, use it for authentication
	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non-ok http status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
