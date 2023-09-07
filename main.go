// main.go
// This file is the entry point and will handle command-line arguments and call the appropriate functions based on those arguments.
package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define command-line flags
	user := flag.String("user", "", "GitHub username to fetch PRs and repos for")
	repo := flag.String("repo", "", "GitHub repo to check pipeline and PR status")
	numPrs := flag.Int("numPrs", 10, "Number of latest PRs to fetch")
	action := flag.String("action", "", "Action to perform (options: userInfo, pipelineStatus, prStatus)")

	// Parse the flags
	flag.Parse()
	// Perform actions based on the flags
	switch *action {
	// User latest Pull requests.
	case "userLatestPRs": // Usage:  Scenario with something: ./chall-scalabit -action=userLatestPRs -user=ucwong -numPrs=1
		// Scenario with nothing (example):./chall-scalabit -action=userLatestPRs -user=Chackoya -numPrs=1
		if *user == "" {
			fmt.Println("Error: GitHub username is required for userLatestPRs action. Usage example: ./chall-scalabit -action=userLatestPRs -user=<username> -numPrs=<number>")
			return
		}
		prs, err := getLatestPRsByUser(*user, *numPrs)
		if err != nil {
			fmt.Printf("Error fetching user's PRs: %s\n", err)
			return
		}
		if len(prs) == 0 {
			fmt.Printf("No PRs to shows. \n")
		}
		for _, pr := range prs {
			fmt.Printf("PR Title: %s, URL: %s, Created At: %s\n", pr.Title, pr.HTMLURL, pr.CreatedAt)
		}
	// Information about the repos of a user case;
	case "userInfoRepos":
		if *user == "" { // ./chall-scalabit -action=userInfoRepos -user=ucwong
			fmt.Println("Error: GitHub username is required for userInfoRepos action.")
			return
		}
		repos, err := getUserRepos(*user)
		if err != nil {
			fmt.Printf("Error fetching user repositories: %s\n", err)
			return
		}

		for _, repo := range repos {
			fmt.Printf("Repo Name: %s, URL: %s\n", repo.Name, repo.URL)
		}
	// Pipeline status (if it exists)
	case "pipelineStatus": // usage on cli : ./chall-scalabit -action=pipelineStatus -repo=actions/starter-workflows ; other example : ./chall-scalabit -action=pipelineStatus -repo=Chackoya/05
		if *repo == "" {
			fmt.Println("Error: GitHub repo is required for pipelineStatus action.")
			return
		}
		checkPipeline(*repo)
	// Pr status case;
	case "prStatus":
		if *repo == "" {
			fmt.Println("Error: GitHub repo is required for prStatus action.")
			return
		}
		getPullRequestStatus(*repo)

	default:
		fmt.Println("Error: Invalid action. Options are userLatestPRs, userInfoRepos, pipelineStatus, prStatus.")
	}

}
