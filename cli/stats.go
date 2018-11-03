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

func stats(c *cli.Context) error {
	fmt.Println("Get pull request statistics")
	config := parseDataArgs(c)
	fmt.Printf("Config: %+v\n", config)
	if err := config.Validate(); err != nil {
		fmt.Printf("Missing arguments: %s\n", err.Error())
		return err
	}

	controllers.GetStats(config)

	return nil
}
