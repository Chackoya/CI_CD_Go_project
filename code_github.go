// code_github.go

/*


Requirements:
-> fetch X latest PRs for a given github user ;  OK

-> fetch the list of repos from a user; OK.

-> check if pipeline is present in a repo or not; OK.

-> given a repo, give the list of open PRs and if they are green , not run , red ;

Recap from git api.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func makeGETRequest(reqURL string) ([]byte, error) {
	// Parse the URL to extract its components
	parsedURL, err := url.Parse(reqURL)
	if err != nil {
		return nil, err
	}

	// Ensure the URL is pointing to api.github.com
	if parsedURL.Host != "api.github.com" {
		return nil, errors.New("invalid URL: requests are only allowed to api.github.com")
	}

	resp, err := http.Get(reqURL) // #nosec
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-ok http status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// getLatestPRsByUser fetches the latest X PRs by a GitHub user across all repositories.
func getLatestPRsByUser(username string, numPrs int) ([]PullRequest, error) {
	// Create the API URL for fetching PRs
	prURL := fmt.Sprintf("https://api.github.com/search/issues?q=author:%s+type:pr&sort=created&order=desc&per_page=%d", username, numPrs)

	// Make the HTTP request
	prData, err := makeGETRequest(prURL)
	if err != nil {
		log.Printf("Error fetching PRs: %s\n", err)
		return nil, fmt.Errorf("error fetching PRs: %s", err)
	}

	// Print the raw JSON data (optional, for debugging)
	//fmt.Println(string(prData), " \n ")

	// Unmarshal the JSON data into our SearchResult struct
	var result SearchResult
	if err := json.Unmarshal(prData, &result); err != nil {
		log.Printf("Error unmarshalling PR JSON: %s\n", err)
		return nil, fmt.Errorf("error unmarshalling PR JSON: %s", err)
	}
	//fmt.Println(result.Items[0])
	// Iterate over the PRs and print their details
	//for _, pr := range result.Items {
	//	fmt.Printf("PR Title: %s, URL: %s, Created At: %s\n", pr.Title, pr.HTMLURL, pr.CreatedAt)
	//}
	return result.Items, nil
}

// getUserRepos fetches the repositories owned by the GitHub user.
func getUserRepos(username string) ([]Repository, error) {
	repoURL := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	repoData, err := makeGETRequest(repoURL)
	if err != nil {
		log.Printf("Error fetching repositories: %s\n", err)
		return nil, fmt.Errorf("error fetching repositories: %s", err)
	}

	var repos []Repository
	if err := json.Unmarshal(repoData, &repos); err != nil {
		log.Printf("Error unmarshalling repo JSON: %s\n", err)
		return nil, fmt.Errorf("error unmarshalling repo json: %s", err)
	}
	/*
		for _, repo := range repos {
			fmt.Printf("Repo Name: %s, URL: %s\n", repo.Name, repo.URL)
		}
	*/
	return repos, nil
}

// checkPipeline checks whether a given GitHub repository has an associated GitHub Actions pipeline.
func checkPipeline(repo string) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/workflows", repo)
	data, err := makeGETRequest(url)
	if err != nil {
		log.Printf("Error fetching pipeline: %s\n", err)
		return
	}

	//fmt.Println(string(data))
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		log.Printf("Error unmarshalling pipeline JSON: %s\n", err)
		return
	}
	// Result contains something like: {"total_count":0,"workflows":[]}

	// Approach: just check if total_count > 0 to see if a pipeline exists;

	totalCount, ok := result["total_count"].(float64) // cast unmarshelled json val;
	if !ok {
		log.Println("Error reading total_count from the result")
		return
	}
	//fmt.Printf("Total Count: %v\n", totalCount)
	if totalCount == 0 {
		fmt.Printf("No pipeline, total: %v\n", totalCount)
	} else {
		fmt.Printf("Pipeline exists, total: %v\n", totalCount)
	}
}

// PullRequestWithStatus represents a GitHub pull request along with its CI/CD status.
type PullRequestWithStatus struct {
	Title  string `json:"title"`
	URL    string `json:"html_url"`
	Ref    string `json:"merge_commit_sha"`
	Status string // This will be populated manually
}

func getPullRequestStatus(repo string) {
	// Fetch open pull requests
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo)
	data, err := makeGETRequest(url)
	if err != nil {
		log.Printf("Error fetching PR statuses: %s\n", err)
		return
	}

	var prStatuses []PullRequestWithStatus
	if err := json.Unmarshal(data, &prStatuses); err != nil {
		log.Printf("Error unmarshalling PR status JSON: %s\n", err)
		return
	}

	// Loop through PRs and get their CI/CD status
	for _, pr := range prStatuses {
		checkRunsURL := fmt.Sprintf("https://api.github.com/repos/%s/commits/%s/check-runs", repo, pr.Ref)
		checkData, err := makeGETRequest(checkRunsURL)
		if err != nil {
			log.Printf("Error fetching check runs: %s\n", err)
			continue
		}

		fmt.Println(string(checkData))
		var checkRuns map[string]interface{}
		if err := json.Unmarshal(checkData, &checkRuns); err != nil {
			log.Printf("Error unmarshalling check runs JSON: %s\n", err)
			continue
		}

		runs, ok := checkRuns["check_runs"].([]interface{})
		if !ok || len(runs) == 0 {
			pr.Status = "not run"
		} else {
			status, _ := runs[0].(map[string]interface{})["conclusion"].(string)
			if status == "success" {
				pr.Status = "green"
			} else if status == "failure" {
				pr.Status = "red"
			} else {
				pr.Status = "not run"
			}
		}

		fmt.Printf("PR Title: %s, Status: %s, URL: %s\n", pr.Title, pr.Status, pr.URL)
	}
}
