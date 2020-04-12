# Github Pull Request Statistics

This repository is a fork of the [horalstvo/ghs](https://github.com/horalstvo/ghs) repository. Kudos to the original author!

## Overview

This command-line tool gathers pull request details from GitHub and stores as a CSV file. Then, you can analyze the CSV file with spreadsheet software like Google Sheets or Microsoft Excel. 

Features: 
- Gathers the following pull request details for the given time window from the given organization and team: repository name, repository number, first review time,  first approval time, second approval time, merge time, number of changed files, number of additions, number of deletions 
- Review and approval times include only workdays
- Stores pull request details into a CSV file
- Retrieves pull requests in chunks (pagination)
- No statistic computations (they are removed from the original code) 
-  Only `stats` command is supported (`single` and `repo` commands were removed from the original code for simplicity)

To access GitHub, the tool needs API token. The token needs `read:org` and `repo` rights. You can create one in Settings -> Developer Settings -> Personal Access Tokens.

## How does it work

The tool speaks to GitHub using an API token. It accepts an organization and a team as input. It retrieves a list of repositories for that team and recent pull requests for every repository (date range is configurable, by default, last two weeks).

Then it calculates how many working hours passed till the first review (first time somebody found time to have a look),  the first and second approvals, and finally, when a pull request was merged. 

## Run with Docker

1. Clone the repo and navigate to the repo folder. 
2. Build a docker image:
```
docker build . -t ghs
```
3. Run the docker container:
```
docker run -v $(pwd):/out ghs stats --org <org-name> --team <team-name> --api-token <api-token> --start -10 --end 0 --file /out/ghs.csv
```
The resulting CSV data will be stored into the `./ghs.csv` file.

## Run with Go

Note, at first, you should install `Go 1.14`. 

1. Clone the repo and navigate to the repo folder. 
2. Install dependencies: 
```
go get
```
3. Build the tool:
```
go build
```
4. Run the tool:
```
chmod +x ghs
./ghs stats --org <org-name> --team <team-name> --api-token <api-token> --start -10 --end 0 --file ghs.csv
```

## Develop

This project is configured for VS Code. However, you can use other IDEs of your choice. If you go with VS Code, it is recommended to install the following plugins: 
- Go
- Code Spell Checker

When developing on macOs, install Xcode and Xcode Command Line Tools. 

