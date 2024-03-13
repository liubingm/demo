package main

import (
	"fmt"
	"time"

	"github.com/zhaokm/sleep/util"
)

func main() {

	// folderPath := flag.String("o", "C:/demo/", "folder path")
	// flag.Parse()
	// util.FolderIteration_new(*folderPath)

	logResults := util.FetchCloudWathLogs()
	if len(logResults) == 0 {
		fmt.Println("no log found")
		return
	}

	systemPrompt, claudePromptsMessages := util.BuildMQTTAnalysisPrompt(logResults)

	fmt.Println(systemPrompt)
	fmt.Println(claudePromptsMessages)

	result, _ := util.CallClaude3WithRetry(systemPrompt, claudePromptsMessages, 5, 120*time.Second)

	fmt.Println("claude Advice" + result)
}
