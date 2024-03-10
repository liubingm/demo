package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Claude3Request struct {
	Messages      []Message `json:"messages"`
	MaxTokens     int       `json:"max_tokens"`
	SystemPrompt  string    `json:"system"`
	Temperature   float64   `json:"temperature,omitempty"`
	TopP          float64   `json:"top_p,omitempty"`
	TopK          int       `json:"top_k,omitempty"`
	StopSequences []string  `json:"stop_sequences,omitempty"`
	Version       string    `json:"anthropic_version"`
}

// ContentItem 对应于JSON中的 "content" 数组里的对象
type Clade3OutputContentItem struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

// Message 对应于整个JSON对象
type Clade3OutputMessage struct {
	Content      []Clade3OutputContentItem `json:"content"`
	ID           string                    `json:"id"`
	Model        string                    `json:"model"`
	Role         string                    `json:"role"`
	StopReason   string                    `json:"stop_reason"`
	StopSequence *string                   `json:"stop_sequence"` // 使用指针，因为它可以是null
	Type         string                    `json:"type"`
	Usage        Clade3OutputUsageInfo     `json:"usage"`
}

// UsageInfo 对应于JSON中的 "usage" 对象
type Clade3OutputUsageInfo struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func CallClaude3WithRetry(systemPrompt string, messages []Message, maxRetries int, retryDelay time.Duration) (string, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		response, err := CallClaude3(systemPrompt, messages)
		if err == nil {
			return response, nil
		}

		lastErr = err
		log.Printf("Attempt %d failed: %s", attempt+1, err)
		time.Sleep(retryDelay)
	}

	// After all attempts, return the last error
	log.Fatalf("All attempts failed: %s", lastErr)
	return "", nil
}

func CallClaude3(systemPrompt string, messages []Message) (string, error) {

	var payloadBytes []byte

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("error LoadDefaultConfig Bedrock, %v", err)
	}

	svc := bedrockruntime.NewFromConfig(cfg)

	payload := Claude3Request{
		Messages:     messages,
		Temperature:  0.4,
		TopP:         0.9,
		TopK:         35,
		SystemPrompt: systemPrompt,
		MaxTokens:    4096,
		Version:      "bedrock-2023-05-31"}

	payloadBytes, err = json.Marshal(payload)
	if err != nil {
		return "", err
	}

	accept := "*/*"
	contentType := "application/json"
	modelId := "anthropic.claude-3-sonnet-20240229-v1:0"

	resp, _ := svc.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		Accept:      &accept,
		ModelId:     &modelId,
		ContentType: &contentType,
		Body:        payloadBytes,
	})

	var out Clade3OutputMessage

	err = json.Unmarshal(resp.Body, &out)
	if err != nil {
		fmt.Printf("unable to Unmarshal JSON, %v", err)
	}

	return out.Content[0].Text, err
}
