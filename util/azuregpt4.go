package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

// GPT4Response represents the structured format of the GPT-4 response
type GPT4Response struct {
	Answer       string
	FullResponse string
}

func CallOpenAIWithRetry(prompt string, systemprompt string, maxRetries int, retryDelay time.Duration) (GPT4Response, string) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		response, jsonResponse, err := CallOpenAI(prompt, systemprompt)
		if err == nil {
			return response, jsonResponse
		}

		lastErr = err
		log.Printf("Attempt %d failed: %s", attempt+1, err)
		time.Sleep(retryDelay)
	}

	// After all attempts, return the last error
	log.Fatalf("All attempts failed: %s", lastErr)
	return GPT4Response{}, ""
}

// CallOpenAI sends a prompt to OpenAI's GPT-4 and returns the response.
func CallOpenAI(userPrompt string, systemPrompt string) (GPT4Response, string, error) {
	azureOpenAIKey := "" //todo input your key
	modelDeploymentID := ""
	azureOpenAIEndpoint := ""

	if azureOpenAIKey == "" || modelDeploymentID == "" || azureOpenAIEndpoint == "" {
		return GPT4Response{}, "", fmt.Errorf("missing credentials")
	}

	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)
	if err != nil {
		return GPT4Response{}, "", fmt.Errorf("creating client: %w", err)
	}

	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{Content: to.Ptr(systemPrompt)},
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(userPrompt)},
	}

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &modelDeploymentID,
		ResponseFormat: &azopenai.ChatCompletionsJSONResponseFormat{},
		Temperature:    to.Ptr[float32](0.4),
		TopP:           to.Ptr[float32](0.9),
		MaxTokens:      to.Ptr(int32(4096)),
	}, nil)
	if err != nil {
		return GPT4Response{}, "", fmt.Errorf("requesting OpenAI: %w", err)
	}

	var gpt4Response GPT4Response
	for _, choice := range resp.Choices {
		if choice.Message != nil && choice.Message.Content != nil {
			gpt4Response.FullResponse += *choice.Message.Content
			// Assuming the main content of the answer is in the first response
			if gpt4Response.Answer == "" {
				gpt4Response.Answer = *choice.Message.Content
			}
		}
	}

	return gpt4Response, gpt4Response.Answer, nil
}
