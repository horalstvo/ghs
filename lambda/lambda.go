package main

import (
	"fmt"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/horalstvo/ghs/controllers"
	"github.com/horalstvo/ghs/models"
)

func HandleRequest(ctx context.Context, config models.StatsConfig) (string, error) {
	if err := config.Validate(); err != nil {
		fmt.Printf("Missing arguments: %s\n", err.Error())
		return "error", err
	}

	controllers.GetStats(config)

	return "done", nil
}

func main() {
	lambda.Start(HandleRequest)
}
