package models

import (
	"sort"
	"time"
)

type (
	Review struct {
		Author    string
		Submitted time.Time
		Status    string
	}

	PullRequest struct {
		Repo    string
		Number  int
		Author  string
		Created time.Time

		FirstReviewedHrs  int
		FirstApprovedHrs  int
		SecondApprovedHrs int
		MergedHrs         int

		ChangedFiles int
		Additions    int
		Deletions    int
	}

	PullRequestByAuthor []PullRequest
)

func (p PullRequestByAuthor) Less(i, j int) bool {
	return p[i].Author > p[j].Author
}

func (p PullRequestByAuthor) Len() int {
	return len(p)
}

func (p PullRequestByAuthor) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PullRequestByAuthor) Sort() {
	sort.Sort(p)
}
