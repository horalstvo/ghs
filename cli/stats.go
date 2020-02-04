package cli

import (
	"fmt"
	"github.com/urfave/cli"

	"github.com/horalstvo/ghs/controllers"
	"github.com/horalstvo/ghs/models"
)

func parseDataArgs(c *cli.Context) models.StatsConfig {
	config := models.StatsConfig{
		Org:      c.String("org"),
		Team:     c.String("team"),
		ApiToken: c.String("api-token"),
		Start:    c.Int("start"),
		End:      c.Int("end"),
	}
	return config
}

func parseSingleDataArgs(c *cli.Context) models.SingleStatsConfig {
	config := models.SingleStatsConfig{
		Org:      c.String("org"),
		Repo:     c.String("repo"),
		ApiToken: c.String("api-token"),
		PrNumber: c.Int("pr-number"),
	}
	return config
}

func parseRepoConfigDataArgs(c *cli.Context) models.RepoConfig {
	config := models.RepoConfig{
		Org:      c.String("org"),
		Repo:     c.String("repo"),
		ApiToken: c.String("api-token"),
		Start:    c.Int("start"),
		End:      c.Int("end"),
	}
	return config
}

func stats(c *cli.Context) error {
	fmt.Println("Get pull requests statistics")
	config := parseDataArgs(c)
	fmt.Printf("Config: %+v\n", config)
	if err := config.Validate(); err != nil {
		fmt.Printf("Missing arguments: %s\n", err.Error())
		return err
	}

	controllers.GetStats(config)

	return nil
}

func singlePr(c *cli.Context) error {
	fmt.Println("Get one pull request statistics")
	config := parseSingleDataArgs(c)
	fmt.Printf("Config: %+v\n", config)
	if err := config.Validate(); err != nil {
		fmt.Printf("Missing arguments: %s\n", err.Error())
		return err
	}

	controllers.GetSingle(config)

	return nil
}

func singleRepo(c *cli.Context) error {
	fmt.Println("Get one repo statistics")
	config := parseRepoConfigDataArgs(c)
	fmt.Printf("Config: %+v\n", config)
	if err := config.Validate(); err != nil {
		fmt.Printf("Missing arguments: %s\n", err.Error())
		return err
	}

	controllers.GetRepoStats(config)

	return nil
}
