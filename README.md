# Github Pull Request Statistics

## Why

I noticed I wait quite a while on pull request reviews. Also there is at times a lot of back and forth which slows 
down delivery.

I wanted to gather some data to support my intuition and help me make this issue more visible and potentially create 
some KPIs for the team.

## How Does It Work

The program speaks to GitHub using an API token. It accepts an organization and a team as input. It retrieves a list of 
repositories for that team and then recent pull requests for every repository (date range is configurable, by default
 last two weeks).

Then it calculates how many working hours passed till first review (first time somebody found time to have a look), 
how many reviews there were in total on the pull request and how long it took to approve.

## Example of Use

- Checkout the repository.
- Install dependencies `go get`.
- Build and run:
```bash
go build
./ghs stats -org <org-name> -team <team-name> -t <api-token>
```

The token needs `read:org` and `repo` rights. You can create one in Settings -> Developer Settings -> Personal Access
 Tokens.

## Future Improvements

Currently more of an MVP. Future improvements:

- Paging get pull requests when last 30 PRs is not enough for specified range.
- Dockerization
