package controllers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"github.com/horalstvo/ghs/external"
	"github.com/horalstvo/ghs/models"
	"github.com/horalstvo/ghs/util"
	"github.com/logrusorgru/aurora"
)

// GetStats returns pull request data
func GetStats(config models.StatsConfig) {
	ctx := context.Background()

	client := external.GetClient(ctx, config.ApiToken)

	repos := external.GetTeamRepos(ctx, config.Org, config.Team, client)

	prs := getPullRequests(ctx, repos, config, client)

	prs = filterPullRequests(prs, config.Start, config.End)

	fmt.Printf("Number of PRs opened in the interval: %d\n", aurora.Blue(len(prs)))

	pullRequests := getDetails(ctx, prs, config.Org, client)

	fmt.Printf("Writing to CSV file '%v'...\n", config.File)
	file, err := os.Create(config.File)
	util.Check(err)
	defer file.Close()

	writer := bufio.NewWriter(file)

	fmt.Fprintf(writer, "Repo,Number,FirstReviewedHrs,FirstApprovedHrs,SecondApprovedHrs,MergedHrs\n")
	for _, pr := range pullRequests {
		fmt.Fprintf(writer, "%v,%v,%v,%v,%v,%v\n", pr.Repo, pr.Number,
			pr.FirstReviewedHrs, pr.FirstApprovedHrs, pr.SecondApprovedHrs, pr.MergedHrs)
	}

	writer.Flush()
	file.Sync()

	fmt.Println("Writing to CSV file completed")

}

func getPullRequests(ctx context.Context, repos []*github.Repository, config models.StatsConfig, client *github.Client) []*github.PullRequest {

	prsPerRepo := make([][]*github.PullRequest, len(repos))
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(repos))

	for i, repo := range repos {
		go func(i int, repoName string) {
			defer waitGroup.Done()
			prsPerRepo[i] = external.GetPullRequests(ctx, config.Org, repoName, client)
			fmt.Printf("Number of PRs returned for %s: %d\n", repoName, aurora.Blue(len(prsPerRepo[i])))
		}(i, *repo.Name)
	}

	waitGroup.Wait()

	prs := make([]*github.PullRequest, 0)
	for i := range repos {
		prs = append(prs, prsPerRepo[i]...)
	}
	return prs
}

func getDetails(ctx context.Context, prs []*github.PullRequest, org string,
	client *github.Client) []models.PullRequest {

	fmt.Printf("Getting details for %v pull requests...\n", len(prs))

	pullRequests := make([]models.PullRequest, len(prs))

	for i, pr := range prs {
		pullRequests[i] = getPullRequestDetails(ctx, org, pr, client)

		// Throttle number of sequential requests to GitHub API
		if (i+1)%25 == 0 {
			fmt.Printf("%v pull requests were processed\n", i)
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Println("Done")

	sort.Slice(pullRequests, func(i, j int) bool {
		return pullRequests[i].Created.Before(pullRequests[j].Created)
	})

	return pullRequests
}

func getPullRequestDetails(ctx context.Context, org string, pr *github.PullRequest,
	client *github.Client) models.PullRequest {

	reviews := getReviews(ctx, org, *pr.Base.Repo.Name, *pr.Number, client)

	firstReviewedHrs := -1
	firstApprovedHrs := -1
	secondApprovedHrs := -1
	mergedHrs := -1

	if len(reviews) > 0 {
		firstReview := reviews[0]
		approvals := getApprovals(reviews)
		firstReviewedHrs = util.WorkHours(*pr.CreatedAt, firstReview.Submitted)

		if len(approvals) > 0 {
			firstApproval := approvals[0]
			firstApprovedHrs = util.WorkHours(*pr.CreatedAt, firstApproval.Submitted)
		}

		if len(approvals) > 1 {
			secondApproval := approvals[1]
			secondApprovedHrs = util.WorkHours(*pr.CreatedAt, secondApproval.Submitted)
		}
	}

	if pr.MergedAt != nil {
		mergedHrs = util.WorkHours(*pr.CreatedAt, *pr.MergedAt)
	}

	pullRequest := models.PullRequest{
		Repo:              *pr.Base.Repo.Name,
		Number:            *pr.Number,
		Created:           *pr.CreatedAt,
		FirstReviewedHrs:  firstReviewedHrs,
		FirstApprovedHrs:  firstApprovedHrs,
		SecondApprovedHrs: secondApprovedHrs,
		MergedHrs:         mergedHrs,

		ChangedFiles: -1,
		Additions:    -1,
		Deletions:    -1,
	}

	if pr.ChangedFiles != nil {
		pullRequest.ChangedFiles = *pr.ChangedFiles
	}

	if pr.Additions != nil {
		pullRequest.Additions = *pr.Additions
	}

	if pr.Deletions != nil {
		pullRequest.Deletions = *pr.Deletions
	}

	return pullRequest
}

func getColored(hours int, percentile float64) aurora.Value {
	if float64(hours) >= percentile {
		return aurora.Red(hours)
	}
	return aurora.Gray(hours)
}

func getApprovals(reviews []models.Review) []models.Review {
	approvals := make([]models.Review, 0)
	for _, rev := range reviews {
		if rev.Status == "APPROVED" {
			approvals = append(approvals, rev)
		}
	}
	return approvals
}

func getReviews(ctx context.Context, org string, repo string, number int, client *github.Client) []models.Review {
	rawReviews := external.GetReviews(ctx, org, repo, number, client)
	reviews := make([]models.Review, 0)
	for _, rev := range rawReviews {
		if rev.SubmittedAt != nil {
			reviews = append(reviews, models.Review{
				Author:    *rev.User.Login,
				Status:    *rev.State,
				Submitted: *rev.SubmittedAt,
			})
		} else {
			fmt.Printf("Skipping %s:%s:%d - nil for submittedAt.\n%v\n", org, repo, number, rev)
		}
	}
	return reviews
}

func filterPullRequests(prs []*github.PullRequest, startDays int, endDays int) []*github.PullRequest {
	from := time.Now().AddDate(0, 0, startDays)
	to := time.Now().AddDate(0, 0, endDays)
	filteredPrs := filter(prs, func(request *github.PullRequest) bool {
		return request.CreatedAt.After(from) && request.CreatedAt.Before(to)
	})
	return filteredPrs
}

func filter(prs []*github.PullRequest, fn func(*github.PullRequest) bool) []*github.PullRequest {
	filtered := make([]*github.PullRequest, 0)
	for _, pr := range prs {
		if fn(pr) {
			filtered = append(filtered, pr)
		}
	}
	return filtered
}
