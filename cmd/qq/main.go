package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var conversationHistory []openai.ChatCompletionMessage

func main() {
	// Get API key and organization ID from environment variables
	apiKey := os.Getenv("OPENAI_API_KEY")
	orgID := os.Getenv("OPENAI_ORGANIZATION")

	if apiKey == "" {
		fmt.Println("Please set the OPENAI_API_KEY environment variable")
		return
	}

	if orgID == "" {
		fmt.Println("Please set the OPENAI_ORGANIZATION environment variable")
		return
	}

	// Create a configuration with both API key and organization ID
	config := openai.DefaultConfig(apiKey)
	config.OrgID = orgID

	client := openai.NewClientWithConfig(config)

	// Check if a question was provided as a command-line argument
	if len(os.Args) > 1 {
		question := strings.Join(os.Args[1:], " ")
		askQuestion(client, question)
		return
	}

	// If no command-line argument, enter interactive mode
	fmt.Println("Ask away...")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		question := scanner.Text()

		if strings.ToLower(question) == "quit" {
			break
		}

		askQuestion(client, question)
	}
}

func askQuestion(client *openai.Client, question string) {

	// Add the user's question to the conversation history
	conversationHistory = append(conversationHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: question,
	})

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT4oLatest,
			Messages:  conversationHistory,
			MaxTokens: 4096,
		},
	)

	// Add the assistant's response to the conversation history
	if err == nil && len(resp.Choices) > 0 {
		conversationHistory = append(conversationHistory, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: resp.Choices[0].Message.Content,
		})
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
