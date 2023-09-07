// code_github_test.go
// A few Integration tests , to use during pipeline...
package main

import (
	"testing"
)

//ref :https://pkg.go.dev/testing

/*
Implementing some basic Integrations tests (on external API =>  github).
*/

func TestGetUserRepos(t *testing.T) {

	// Empty token to run on Github Actions... not the best approach but at least the tests will run... with unauth reqs..
	token := ""

	// Integration Test: testing the external API call to get repos.
	username := "Chackoya" // using a real username, example of my acc: Chackoya
	repos, err := getUserRepos(token, username)
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
func TestCheckPipeline(t *testing.T) {

	// Same... using unauth reqs to the github API , to run during github actions steps.
	token := ""

	// In this example, usage of  "actions/starter-workflows" , can be replaced with another...
	repo := "actions/starter-workflows"
	totalCount, err := checkPipeline(token, repo)
	if err != nil {
		t.Errorf("Error when checking pipeline for repo %s: %s", repo, err)
		return
	}

	if totalCount <= 0 {
		t.Errorf("Expected repo %s to have GitHub actions, but got totalCount: %d", repo, totalCount)
		return
	}
}
