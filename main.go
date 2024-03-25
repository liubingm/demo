package main

import (
	"fmt"
	"time"

	"github.com/zhaokm/sleep/util"
	//"github.com/zhaokm/sleep/util"
)

func main() {

	//从CloudWatch读取日志
	logResults := util.FetchCloudWatchLogs()
	if len(logResults) == 0 {
		fmt.Println("no log found")
		return
	}
	//将日志数据导入Prompt模板，形成大模型输入
	systemPrompt, claudePromptsMessages := util.BuildLogAnalysisPrompt(logResults)
	//fmt.Println("系统提示" + systemPrompt)
	fmt.Println("Input", claudePromptsMessages)

	//调用Bedrock-Claude3模型，分析并返回结果
	result, _ := util.CallClaude3WithRetry(systemPrompt, claudePromptsMessages, 5, 120*time.Second)
	fmt.Println("claude Advice")
	fmt.Println(result)
}
