# Scalabit Challenge - Programming challenge


This repository contains a Go project that interacts with the GitHub API to perform various operations. This README provides steps on setting up, building, and running the project.


Requirements:
- Golang. 
Version used:
go version go1.18.1 linux/amd64

- Docker


## How to run the project

### Step 1: download project

Clone/download the repo and go into the root folder

### Step 2: modules

Install the modules required:

> go mod tidy


### Step 3: environment

Setup the .env file.
You can find the in project a .env.example file, create .env file with that template. Then write on it your Github Token. 

Note: this step is not mandatory, the token is just to make auth requests to the Github API, in order to avoid limitations from their side. Otherwise the code is adapted to make requests without auth, but the result is not great as we will reach the maximum requests very quickly for some operations (such as getting the PR status...).


### Step 4: running the project

You can use for example the following command (which will build an executable...):

> go build -o chall-scalabit

After the executable is created you can try out the project with the following commands (those are real examples, please replace the parameters, such as user, with your values):

#### To find the latests X PRs: 

> ./chall-scalabit -action=userLatestPRs -user=ucwong -numPrs=5

#### List of repos for a user:

> ./chall-scalabit -action=userInfoRepos -user=ucwong

#### Check if a repo has a pipeline:

> ./chall-scalabit -action=pipelineStatus -repo=actions/starter-workflows

#### PR Status:

> ./chall-scalabit -action=prStatus -repo=ethereum/go-ethereum



#### Run tests

There are a few integrations tests, mainly to be run on the pipeline.

> go test



##





















