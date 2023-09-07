// code_github.go

/*
Requirements:
-> fetch X latest PRs for a given github user ;
OK

-> fetch the list of repos from a user;
OK.

-> check if pipeline is present in a repo or not;
OK.

-> given a repo, give the list of open PRs and if they are green , not run , red ;
OK


Recap from git api.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

// PullRequest represents a GitHub pull request with essential fields.
type PullRequest struct {
	Title     string `json:"title"`
	HTMLURL   string `json:"html_url"`
	CreatedAt string `json:"created_at"`
}

// SearchResult represents the result returned from a GitHub search query for pull requests.
type SearchResult struct {
	Items []PullRequest `json:"items"`
}

// Repository represents a GitHub repository with its name and URL.
type Repository struct {
	Name string `json:"name"`
	URL  string `json:"html_url"`
}

// getLatestPRsByUser fetches the latest X PRs by a GitHub user across all repositories.
func getLatestPRsByUser(token, username string, numPrs int) ([]PullRequest, error) {
	// Create the API URL for fetching PRs
	prURL := fmt.Sprintf("https://api.github.com/search/issues?q=author:%s+type:pr&sort=created&order=desc&per_page=%d", username, numPrs)

	// Make the HTTP request
	prData, err := makeGETRequest(prURL, token)
	if err != nil {
		log.Printf("Error fetching PRs: %s\n", err)
		return nil, fmt.Errorf("error fetching PRs: %s", err)
	}
	//fmt.Println(string(prData), " \n ")

	// Unmarshal the JSON data into our SearchResult struct
	var result SearchResult
	if err := json.Unmarshal(prData, &result); err != nil {
		log.Printf("Error unmarshalling PR JSON: %s\n", err)
		return nil, fmt.Errorf("error unmarshalling PR JSON: %s", err)
	}
	// Return the result (prs) and print them on the console (in main.go)
	return result.Items, nil
}

// getUserRepos fetches the repositories owned by the GitHub user.
func getUserRepos(token, username string) ([]Repository, error) {
	repoURL := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	repoData, err := makeGETRequest(repoURL, token)
	if err != nil {
		log.Printf("Error fetching repositories: %s\n", err)
		return nil, fmt.Errorf("error fetching repositories: %s", err)
	}

	var repos []Repository
	if err := json.Unmarshal(repoData, &repos); err != nil {
		log.Printf("Error unmarshalling repo JSON: %s\n", err)
		return nil, fmt.Errorf("error unmarshalling repo json: %s", err)
	}
	// return repos, to be printed on the console (on the main.go)
	return repos, nil
}

// checkPipeline checks whether a given GitHub repository has an associated GitHub Actions pipeline.
func checkPipeline(token, repo string) (int, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/workflows", repo)
	data, err := makeGETRequest(url, token)
	if err != nil {
		log.Printf("Error fetching pipeline: %s\n", err)
		return 0, fmt.Errorf("error fetching pipeline: %s", err)
	}

	//fmt.Println(string(data))
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		log.Printf("Error unmarshalling pipeline JSON: %s\n", err)
		return 0, fmt.Errorf("error unmarshalling pipeline JSON: %s", err)
	}
	// Result contains something like: {"total_count":0,"workflows":[]}

	// Approach: just check if total_count > 0 to see if a pipeline exists;

	totalCount, ok := result["total_count"].(float64) // cast unmarshelled json val;
	if !ok {
		log.Println("Error reading total_count from the result")
		return 0, errors.New("error reading total_count from the result")
	}

	//fmt.Printf("Total Count: %v\n", totalCount)
	/*
		if totalCount == 0 {
			fmt.Printf("No pipeline, total: %v\n", totalCount)
		} else {
			fmt.Printf("Pipeline exists, total: %v\n", totalCount)
		}
	*/
	return int(totalCount), nil
}

///////////////////////////////////////////////////////////////////////////////////////////////

/*
Last endpoint

- Fetch PULL REQUEST STATUS (green, red , not run)
*/
// PullRequestWithStatus represents a GitHub pull request along with its CI/CD status.
type PullRequestWithStatus struct {
	Title  string `json:"title"`
	URL    string `json:"html_url"`
	Ref    string `json:"merge_commit_sha"`
	Status string
}

func getCombinedCommitStatus(token, repo, ref string) string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/commits/%s/status", repo, ref)

	data, err := makeGETRequest(url, token)
	if err != nil {
		log.Printf("Error 2 fetching combined commit status: %s\n", err)
		return "error"
	}

	var statusResponse map[string]interface{}
	err = json.Unmarshal(data, &statusResponse)
	if err != nil {
		log.Printf("Error unmarshalling commit status: %s\n", err)
		return "error"
	}
	//fmt.Println(statusResponse["state"].(string))
	return statusResponse["state"].(string)
}

func getPullRequestStatus(token, repo string) {
	//url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo)

	data, err := makeGETRequest(url, token)
	if err != nil {
		log.Printf("Error 1 fetching PR statuses: %s\n", err)
		return
	}

	var prStatuses []PullRequestWithStatus
	err = json.Unmarshal(data, &prStatuses)
	if err != nil {
		log.Printf("Error unmarshalling PR status JSON: %s\n", err)
		return
	}

	for _, pr := range prStatuses {

		status := getCombinedCommitStatus(token, repo, pr.Ref)
		fmt.Println("STATUS:", status)
		switch status {
		case "success":
			pr.Status = "green"
		case "failure":
			pr.Status = "red"
		case "pending":
			pr.Status = "not run"
		default:
			pr.Status = "unknown"
		}

		fmt.Printf("PR Title: %s, Status: %s, URL: %s\n", pr.Title, pr.Status, pr.URL)
	}
}
