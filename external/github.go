package external

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/horalstvo/ghs/util"
	"golang.org/x/oauth2"
)

func GetClient(ctx context.Context, apiToken string) *github.Client {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func GetPullRequests(org string, repo string, ctx context.Context, client *github.Client) []*github.PullRequest {
	allPrs := make([]*github.PullRequest, 0)
	pageSize := 50
	page := -1
	for {
		page += 1
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

func GetReviews(org string, repo string, number int, ctx context.Context,
	client *github.Client) []*github.PullRequestReview {
	reviews, _, err := client.PullRequests.ListReviews(ctx, org, repo, number, &github.ListOptions{})
	util.Check(err)
	return reviews
}

func GetTeamRepos(org string, team string, ctx context.Context, client *github.Client) []*github.Repository {
	teamId, getTeamErr := getTeamId(org, team, ctx, client)
	util.Check(getTeamErr)

	repos, _, err := client.Teams.ListTeamRepos(ctx, *teamId, &github.ListOptions{})
	util.Check(err)
	return repos
}

func getTeamId(org string, team string, ctx context.Context, client *github.Client) (*int64, error) {
	teams, _, err := client.Teams.ListTeams(ctx, org, &github.ListOptions{})
	util.Check(err)

	for _, t := range teams {
		if *t.Name == team {
			return t.ID, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Team not found for %s", team))
}
