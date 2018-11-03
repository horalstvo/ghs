package models

import "time"

type Review struct {
	Author    string
	Submitted time.Time
	Status    string
}

type PullRequest struct {
	Repo    string
	Number  int
	Author  string
	Created time.Time
	Reviews []Review

	FirstReview   int
	ApprovalAfter int
	ApprovedBy    string
}
