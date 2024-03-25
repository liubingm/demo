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
//const LOG_PATTERN = "info"           // change it to error when online

func FetchCloudWatchLogs() []string {
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

	for _, stream := range result.LogStreams { //遍历所有日志流，对每个日志流调用 `fetchLogStreamLog` 函数，获取该日志流中的日志事件
		logResults = append(logResults, fetchLogStreamLog(svc, LOGGROUP_NAME, *stream.LogStreamName)...)
	}
	return logResults
}

func fetchLogStreamLog(awsConfig *cloudwatchlogs.Client, logGroupName string, logStreamName string) []string {
	var logResults []string
	input := &cloudwatchlogs.GetLogEventsInput{ //`fetchLogStreamLog` 函数调用 `GetLogEvents` 方法，获取指定日志流中的所有日志事件。
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	}

	result, err := awsConfig.GetLogEvents(context.TODO(), input)
	if err != nil {
		log.Fatalf("GetLogEvents error, %v", err)
	}
	for _, event := range result.Events { //遍历所有日志事件，打印每个事件的时间戳和消息内容，并将消息内容添加到一个字符串切片中。
		//fmt.Printf("Timestamp: %d, Message: %s\n", event.Timestamp, *event.Message)
		logResults = append(logResults, aws.ToString(event.Message))
	}
	return logResults
}
