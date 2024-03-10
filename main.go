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

	systemPrompt, claudePromptsMessages := util.BuildMQTTAnalysisPrompt(logResults)

	result, _ := util.CallClaude3WithRetry(systemPrompt, claudePromptsMessages, 5, 120*time.Second)

	fmt.Println(result)

}
