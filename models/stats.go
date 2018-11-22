package models

import "errors"

type StatsConfig struct {
	Org      string
	Team     string
	ApiToken string
	Start    int
	End      int
}

type SingleStatsConfig struct {
	Org      string
	Repo     string
	ApiToken string
	PrNumber int
}

func (c StatsConfig) Validate() error {
	if len(c.ApiToken) == 0 {
		return errors.New("api token is missing")
	}
	if len(c.Org) == 0 {
		return errors.New("organization is missing")
	}
	if len(c.Team) == 0 {
		return errors.New("team is missing")
	}
	if c.Start >= 0 {
		return errors.New("start is invalid")
	}
	if c.End > 0 {
		return errors.New("end is invalid")
	}
	return nil
}

func (c SingleStatsConfig) Validate() error {
	if len(c.ApiToken) == 0 {
		return errors.New("api token is missing")
	}
	if len(c.Org) == 0 {
		return errors.New("organization is missing")
	}
	if len(c.Repo) == 0 {
		return errors.New("repo is missing")
	}
	if c.PrNumber <= 0 {
		return errors.New("PR number is invalid")
	}
	return nil
}
