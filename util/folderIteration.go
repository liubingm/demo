package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func FolderIteration(folderPath string) string {
	var CodeBuilder strings.Builder
	var codes []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		code, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		fmt.Printf("Source Code %s\n", path)
		CodeBuilder.WriteString(string(code))

		codes = append(codes, string(code))

		return nil
	})

	if err != nil {
		fmt.Println("遍历文件夹时出错：", err)
	}

	systemPrompt, userPrompt := BuildGPTPrompt(codes)

	fmt.Println("prompt design")
	fmt.Println(systemPrompt)
	fmt.Println(userPrompt)

	_, jsonResponse, _ := CallOpenAI(userPrompt, systemPrompt)

	fmt.Println("gpt4 response")
	fmt.Println(jsonResponse)

	// 返回Java代码的字符串
	return ""
}

func FolderIteration_new(folderPath string) string {
	var codeBuilder strings.Builder
	var codes []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只处理.go文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			code, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Printf("Go Source Code: %s\n", path)
			codeContent := string(code)
			// 将文件内容添加到字符串构建器和代码列表中
			codeBuilder.WriteString(codeContent + "\n\n") // 加上换行，以便区分不同文件的内容
			codes = append(codes, codeContent)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("遍历文件夹时出错：%v\n", err)
	}

	systemPrompt, claudePromptsMessages := BuildCodeAnalysisPrompt(codes)

	result, _ := CallClaude3WithRetry(systemPrompt, claudePromptsMessages, 5, 120*time.Second)

	fmt.Println(result)

	// 返回合并后的Go代码字符串
	return codeBuilder.String()
}
