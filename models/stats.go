package models

import "errors"

type StatsConfig struct {
	Org      string
	Team     string
	ApiToken string
	Start    int
	End      int
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
