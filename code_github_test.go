// code_github_test.go

package main

import (
	"testing"
)

//ref :https://pkg.go.dev/testing

/*
Implementing some basic Integrations tests (on external API =>  github).
*/
func TestGetUserRepos(t *testing.T) {
	// Integration Test: testing the external API call to get repos.
	username := "Chackoya" // using a real username, example of my acc: Chackoya
	repos, err := getUserRepos(username)
	if err != nil {
		t.Errorf("Error should be nil, got: %s", err)
	}

	// t.Errorf reports an error , but non blocking... if it catches an error other tests keep running

	// Check for at least one repo
	if len(repos) == 0 {
		t.Errorf("Expected at least one repo, got zero")
	}

	//fmt.Println(repos)
	// Check if a known field exists in the first repo
	if repos[0].Name == "" || repos[0].URL == "" {
		t.Errorf("Expected non-empty repo name and URL")
	}
}
