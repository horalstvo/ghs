package controllers

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/google/go-github/github"
	"github.com/horalstvo/ghs/external"
	"github.com/horalstvo/ghs/models"
	"github.com/horalstvo/ghs/util"
	"github.com/logrusorgru/aurora"
	"github.com/montanaflynn/stats"
)

func GetStats(config models.StatsConfig) {
	ctx := context.Background()

	client := external.GetClient(ctx, config.ApiToken)

	repos := external.GetTeamRepos(config.Org, config.Team, ctx, client)

	prs := getPullRequests(repos, config, ctx, client)

	filteredPrs := filterPullRequests(prs, config.Start, config.End)

	fmt.Printf("Number of PRs opened in the interval: %d\n", aurora.Blue(len(filteredPrs)))

	pullRequests := getDetails(filteredPrs, config.Org, ctx, client)

	fmt.Printf("Calculate statistics\n")

	first80, approval80 := getStatistics(pullRequests)
	models.PullRequestByAuthor(pullRequests).Sort()

	fmt.Printf("80th percentiles: first review: %f approval: %f\n", first80, approval80)
	tw := new(tabwriter.Writer)
	tw.Init(os.Stdout, 6, 4, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)
	fmt.Fprintf(tw, "Repo\tPR\tAuthor\t%s\t#\t%s\tApproved\n", aurora.Black("1st"), aurora.Red("App"))
	for _, pr := range pullRequests {
		fmt.Fprintf(tw, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pr.Repo, pr.Number, pr.Author,
			getColored(pr.FirstReview, first80), len(pr.Reviews),
			getColored(pr.ApprovalAfter, approval80), pr.ApprovedBy)
	}
	tw.Flush()
}

func GetSingle(config models.SingleStatsConfig) {
	ctx := context.Background()

	client := external.GetClient(ctx, config.ApiToken)

	pr, _, err := client.PullRequests.Get(ctx, config.Org, config.Repo, config.PrNumber)
	util.Check(err)

	prs := []*github.PullRequest{pr}

	pullRequests := getDetails(prs, config.Org, ctx, client)

	tw := new(tabwriter.Writer)
	tw.Init(os.Stdout, 6, 4, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)
	fmt.Fprintf(tw, "Repo\tPR\tAuthor\t%s\t#\t%s\tApproved\n", aurora.Black("1st"), aurora.Red("App"))
	for _, pr := range pullRequests {
		fmt.Fprintf(tw, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pr.Repo, pr.Number, pr.Author,
			pr.FirstReview, len(pr.Reviews), pr.ApprovalAfter, pr.ApprovedBy)
	}
	tw.Flush()
}

func GetRepoStats(config models.RepoConfig) {
	ctx := context.Background()

	client := external.GetClient(ctx, config.ApiToken)

	prs := external.GetPullRequests(config.Org, config.Repo, ctx, client)

	filteredPrs := filterPullRequests(prs, config.Start, config.End)

	fmt.Printf("Number of PRs opened in the interval: %d\n", aurora.Blue(len(filteredPrs)))

	pullRequests := getDetails(filteredPrs, config.Org, ctx, client)

	fmt.Printf("Calculate statistics\n")

	first80, approval80 := getStatistics(pullRequests)
	models.PullRequestByAuthor(pullRequests).Sort()

	fmt.Printf("80th percentiles: first review: %f approval: %f\n", first80, approval80)
	tw := new(tabwriter.Writer)
	tw.Init(os.Stdout, 6, 4, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)
	fmt.Fprintf(tw, "Repo\tPR\tAuthor\t%s\t#\t%s\tApproved\n", aurora.Black("1st"), aurora.Red("App"))
	for _, pr := range pullRequests {
		fmt.Fprintf(tw, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pr.Repo, pr.Number, pr.Author,
			getColored(pr.FirstReview, first80), len(pr.Reviews),
			getColored(pr.ApprovalAfter, approval80), pr.ApprovedBy)
	}
	tw.Flush()
}



func getPullRequests(repos []*github.Repository, config models.StatsConfig, ctx context.Context, client *github.Client) []*github.PullRequest {

	prsPerRepo := make([][]*github.PullRequest, len(repos))
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(repos))

	for i, repo := range repos {
		go func(i int, repoName string) {
			defer waitGroup.Done()
			prsPerRepo[i] = external.GetPullRequests(config.Org, repoName, ctx, client)
			fmt.Printf("Number of PRs returned for %s: %d\n", repoName, aurora.Blue(len(prsPerRepo[i])))
		}(i, *repo.Name)
	}

	waitGroup.Wait()

	prs := make([]*github.PullRequest, 0)
	for i, _ := range repos {
		prs = append(prs, prsPerRepo[i]...)
	}
	return prs
}

func getDetails(lastWeekPrs []*github.PullRequest, org string, ctx context.Context,
	client *github.Client) []models.PullRequest {

	fmt.Println("Getting pull requests details.")

	pullRequests := make([]models.PullRequest, len(lastWeekPrs))
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(lastWeekPrs))

	for i, pr := range lastWeekPrs {
		go func(i int, pr *github.PullRequest) {
			defer waitGroup.Done()
			pullRequests[i] = getPullRequestDetails(org, pr, ctx, client)
		}(i, pr)
	}
	waitGroup.Wait()
	fmt.Println("Done")

	sort.Slice(pullRequests, func(i, j int) bool {
		return pullRequests[i].Author < pullRequests[j].Author
	})

	return pullRequests
}

func getPullRequestDetails(org string, pr *github.PullRequest, ctx context.Context,
	client *github.Client) models.PullRequest {

	reviews := getReviews(org, *pr.Base.Repo.Name, *pr.Number, ctx, client)

	if len(reviews) > 0 {
		approval := getApproval(reviews)
		firstReview := reviews[0]
		nilReview := models.Review{}
		firstReviewHours := util.WorkHours(*pr.CreatedAt, firstReview.Submitted)

		if approval != nilReview {
			approvalHours := util.WorkHours(*pr.CreatedAt, approval.Submitted)
			return models.PullRequest{
				Repo:          *pr.Base.Repo.Name,
				Number:        *pr.Number,
				Author:        *pr.User.Login,
				Created:       *pr.CreatedAt,
				Reviews:       reviews,
				FirstReview:   firstReviewHours,
				ApprovalAfter: approvalHours,
				ApprovedBy:    approval.Author,
			}
		} else {
			return models.PullRequest{
				Repo:        *pr.Base.Repo.Name,
				Number:      *pr.Number,
				Author:      *pr.User.Login,
				Created:     *pr.CreatedAt,
				Reviews:     reviews,
				FirstReview: firstReviewHours,
			}
		}
	} else {
		return models.PullRequest{
			Repo:    *pr.Base.Repo.Name,
			Number:  *pr.Number,
			Author:  *pr.User.Login,
			Created: *pr.CreatedAt,
		}
	}
}

func getStatistics(pullRequests []models.PullRequest) (float64, float64) {
	firstReviews := make([]float64, 0)
	approvals := make([]float64, 0)
	for _, pr := range pullRequests {
		if len(pr.Reviews) > 0 {
			firstReviews = append(firstReviews, float64(pr.FirstReview))
		}
		if len(pr.ApprovedBy) > 0 {
			approvals = append(approvals, float64(pr.ApprovalAfter))
		}
	}
	first80, err := stats.Percentile(firstReviews, 80)
	util.Check(err)
	approval80, errApr := stats.Percentile(approvals, 80)
	util.Check(errApr)
	return first80, approval80
}

func getColored(hours int, percentile float64) aurora.Value {
	if float64(hours) >= percentile {
		return aurora.Red(hours)
	}
	return aurora.Gray(hours)
}

func getApproval(reviews []models.Review) models.Review {
	for _, rev := range reviews {
		if rev.Status == "APPROVED" {
			return rev
		}
	}
	return models.Review{}
}

func getReviews(org string, repo string, number int, ctx context.Context, client *github.Client) []models.Review {
	pr_reviews := external.GetReviews(org, repo, number, ctx, client)
	reviews := make([]models.Review, 0)
	for _, rev := range pr_reviews {
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
	lastWeekPrs := filter(prs, func(request *github.PullRequest) bool {
		return request.CreatedAt.After(from) && request.CreatedAt.Before(to)
	})
	return lastWeekPrs
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
