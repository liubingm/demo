package util

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

const REGION = "us-east-1"
const LOGGROUP_NAME = "/aws/lambda/SSMOnboardingLambda" // change the log group name when online
const LOG_PATTERN = "INFO"                              // change it to error when online

func FetchCloudWathLogs() []string {
	var logResults []string
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(REGION),
	)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return logResults
	}

	svc := cloudwatchlogs.NewFromConfig(cfg)

	startTime := time.Now().Add(-24 * time.Hour).UnixMilli() // COMMENT THIS AFTER ONLINE
	// startTime := time.Now().AddDate(0, -6, 0).UnixMilli() // UNCOMMENT THIS
	endTime := time.Now().UnixMilli()

	input := &cloudwatchlogs.FilterLogEventsInput{
		StartTime:     aws.Int64(startTime),
		EndTime:       aws.Int64(endTime),
		LogGroupName:  aws.String(LOGGROUP_NAME),
		FilterPattern: aws.String(LOG_PATTERN),
	}

	result, err := svc.FilterLogEvents(context.TODO(), input)
	if err != nil {
		fmt.Println("Error fetching logs:", err)
		return logResults
	}

	for _, event := range result.Events {
		logResults = append(logResults, aws.ToString(event.Message))
	}
	return logResults
}
