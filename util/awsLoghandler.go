package util

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

const REGION = "us-east-1"

// const LOGGROUP_NAME = "/aws/lambda/SSMOnboardingLambda" // change the log group name when online
const LOGGROUP_NAME = "AWSIotLogsV2" // change the log group name when online
const LOG_PATTERN = "info"           // change it to error when online

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

	// startTime := time.Now().Add(-24 * time.Hour).UnixMilli() // COMMENT THIS AFTER ONLINE

	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(LOGGROUP_NAME),
		OrderBy:      types.OrderBy("LastEventTime"),
		Descending:   aws.Bool(true),
	}

	result, err := svc.DescribeLogStreams(context.TODO(), input)
	if err != nil {
		log.Fatalf("DescribeLogStreams error, %v", err)
	}

	// 输出日志流信息
	for _, stream := range result.LogStreams {
		logResults = append(logResults, fetchLogStreamLog(svc, LOGGROUP_NAME, *stream.LogStreamName)...)
	}
	return logResults
}

func fetchLogStreamLog(awsConfig *cloudwatchlogs.Client, logGroupName string, logStreamName string) []string {
	var logResults []string
	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	}

	result, err := awsConfig.GetLogEvents(context.TODO(), input)
	if err != nil {
		log.Fatalf("GetLogEvents error, %v", err)
	}
	for _, event := range result.Events {
		fmt.Printf("Timestamp: %d, Message: %s\n", event.Timestamp, *event.Message)
		logResults = append(logResults, aws.ToString(event.Message))
	}
	return logResults
}
