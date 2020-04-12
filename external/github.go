package external

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/horalstvo/ghs/util"
	"golang.org/x/oauth2"
)

// GetClient returns GitHub client
func GetClient(ctx context.Context, apiToken string) *github.Client {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

// GetPullRequests returns pull requests of the given repository
func GetPullRequests(ctx context.Context, org string, repo string, client *github.Client) []*github.PullRequest {
	allPrs := make([]*github.PullRequest, 0)
	pageSize := 50
	page := -1
	for {
		page++
		prs, _, err := client.PullRequests.List(ctx, org, repo, &github.PullRequestListOptions{
			Sort:      "created",
			State:     "all",
			Direction: "desc",
			ListOptions: github.ListOptions{
				PerPage: pageSize,
				Page:    page,
			},
		})
		util.Check(err)
		allPrs = append(allPrs, prs...)

		if len(prs) < pageSize {
			break
		}
	}
	return allPrs
}

// GetReviews returns reviews of the given pull request
func GetReviews(ctx context.Context, org string, repo string, number int,
	client *github.Client) []*github.PullRequestReview {
	reviews, _, err := client.PullRequests.ListReviews(ctx, org, repo, number, &github.ListOptions{})
	util.Check(err)
	return reviews
}

// GetPullRequest returns details of the given pull request
func GetPullRequest(ctx context.Context, org string, repo string, number int,
	client *github.Client) *github.PullRequest {
	pullRequest, _, err := client.PullRequests.Get(ctx, org, repo, number)
	util.Check(err)
	return pullRequest
}

// GetTeamRepos returns repositories of the given team
func GetTeamRepos(ctx context.Context, org string, team string, client *github.Client) []*github.Repository {
	teamID, getTeamErr := getTeamID(ctx, org, team, client)
	util.Check(getTeamErr)

	repos, _, err := client.Teams.ListTeamRepos(ctx, *teamID, &github.ListOptions{})
	util.Check(err)
	return repos
}

// GetTeamID returns the identifier of the team by name
func getTeamID(ctx context.Context, org string, team string, client *github.Client) (*int64, error) {
	teams, _, err := client.Teams.ListTeams(ctx, org, &github.ListOptions{})
	util.Check(err)

	for _, t := range teams {
		if *t.Name == team {
			return t.ID, nil
		}
	}
	return nil, fmt.Errorf("Team not found for %s", team)
}
