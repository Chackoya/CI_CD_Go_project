// Gustavo Gama - SCALABIT TEST.
// main.go
// This file is the entry point and will handle command-line arguments and call the appropriate functions based on those arguments.
// please check out the README.md file for information about the project.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Define command-line flags
	user := flag.String("user", "", "GitHub username to fetch PRs and repos for")
	repo := flag.String("repo", "", "GitHub repo to check pipeline and PR status")
	numPrs := flag.Int("numPrs", 10, "Number of latest PRs to fetch")
	action := flag.String("action", "", "Action to perform (options: userInfo, pipelineStatus, prStatus)")

	// Basically it's preferable to define the token on the .env file so avoid limitations from the Github API. The following code will try to load it from there.
	// If the constant "GITHUB_TOKEN" is empty (i.e: GITHUB_TOKEN=), then we will use the free API (with limitations 60 i guess).
	// If we insert the token we will proceed with it for more requests (more availability from the external api).
	err := godotenv.Load() // Try to load the .env file if it exists
	if err != nil {
		log.Println("Warning: Could not load .env file. Relying on environment variables.")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Println("Warning: GITHUB_TOKEN is not set. Making unauthenticated requests.")
	}

	/*
		if token == "" { // this fct is too heavy on req, better not use it without auth.
			log.Fatal("GITHUB_TOKEN is not set.")
		}
	*/

	// Parse the flags
	flag.Parse()
	// Perform actions based on the flags
	switch *action {
	// User latest Pull requests.
	case "userLatestPRs": // Usage:  Scenario with something: ./chall-scalabit -action=userLatestPRs -user=ucwong -numPrs=5
		// Scenario with nothing (example):./chall-scalabit -action=userLatestPRs -user=Chackoya -numPrs=1
		if *user == "" {
			fmt.Println("Error: GitHub username is required for userLatestPRs action. Usage example: ./chall-scalabit -action=userLatestPRs -user=<username> -numPrs=<number>")
			return
		}
		prs, err := getLatestPRsByUser(token, *user, *numPrs)
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
		repos, err := getUserRepos(token, *user)
		if err != nil {
			fmt.Printf("Error fetching user repositories: %s\n", err)
			return
		}

		for _, repo := range repos {
			fmt.Printf("Repo Name: %s, URL: %s\n", repo.Name, repo.URL)
		}
	// Pipeline status (if it exists)
	case "pipelineStatus": // usage on cli : ./chall-scalabit -action=pipelineStatus -repo=actions/starter-workflows ;
		// other example total_count:0 => ./chall-scalabit -action=pipelineStatus -repo=Chackoya/05
		if *repo == "" {
			fmt.Println("Error: GitHub repo is required for pipelineStatus action.")
			return
		}

		totalCount, err := checkPipeline(token, *repo)
		if err != nil {
			fmt.Printf("Error when checking pipeline: %s\n", err)
			return
		}

		if totalCount > 0 {
			fmt.Printf("Pipeline exists for repo: %s , total counter: %v \n", *repo, totalCount)
		} else {
			fmt.Printf("No pipeline found for repo: %s , total counter: %v \n", *repo, totalCount)
		}
	// Pr status case;
	case "prStatus": // usage on cli : ./chall-scalabit -action=prStatus -repo=actions/starter-workflows
		// Other example (that contains errors at time of writing 7september): ./chall-scalabit -action=prStatus -repo=ethereum/go-ethereum
		if *repo == "" {
			fmt.Println("Error: GitHub repo is required for prStatus action.")
			return
		}
		getPullRequestStatus(token, *repo)

	default:
		fmt.Println("Error: Invalid action. Options are userLatestPRs, userInfoRepos, pipelineStatus, prStatus.")
	}

}
