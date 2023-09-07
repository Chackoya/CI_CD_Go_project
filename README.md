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



Note: on the CLI, 

-action=<string> option corresponds to the task we want (e.g. fetch latest X PRs etc);

-user=<string> is the name of the github user; 

-numPRs=<int> (limit of PRs we want);

-repo=<str>  corresponds to the repository (it follows the pattern owner/repo_name)


#### Run tests

There are a few integrations tests, mainly to be run on the pipeline.

> go test



## Pipeline information

The pipeline provides CI/CD for the project. It runs checks, builds the project, runs tests and deploys (with docker). The triggers are based on push/pull req.

Please check out the .yml file for more information. 

#### Run the pipeline

You're going to setup the following secrets (Secrets section on Github):

- DOCKER_HUB_USERNAME: Your Docker Hub username.

- DOCKER_HUB_ACCESS_TOKEN: Your Docker Hub access token.


Once this is done, just make a push and the it should trigger the process. 

Example to trigger (config the git to be able to push, git init or other...):

> git add .

> git commit -m "Triggering CI/CD"

> git push


#### Check the results

Go into the Actions tab and check the results of the pipeline process.

If all went well, you can run the docker image with the following commands:

Open a terminal and pull it:

> docker pull <DOCKER_HUB_USERNAME>/my-app:latest

Run command (you can provide your github token to make auth requests)...

> docker run -e GITHUB_TOKEN=<Your_GitHub_Token> <DOCKER_HUB_USERNAME>/my-app:latest





#### You can also test my own docker image (on DockerHub)


Pull my image:

> docker pull chackoya/my-app:latest



Example (fetch repos) on how to run the CLI: 

> docker run -e GITHUB_TOKEN=<Your_GitHub_Token> chackoya/my-app:latest ./main -action=userInfoRepos -user=Chackoya

Please modify this command as you need (-actions) and other parameters to call other functions.









## Final comments

Thanks for reading and testing out the project.

In case of any doubt or if you find any problem with the app (bug , issue with running, testing the project). Please feel free to contact me: 

gustavo.p.gama@gmail.com
















