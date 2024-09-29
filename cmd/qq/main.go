package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

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
	fmt.Println("Welcome to the ChatGPT CLI tool!")
	fmt.Println("Type your questions and press Enter. Type 'quit' to exit.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		question := scanner.Text()

		if strings.ToLower(question) == "quit" {
			break
		}

		askQuestion(client, question)
	}

	fmt.Println("Thank you for using the ChatGPT CLI tool. Goodbye!")
}

func askQuestion(client *openai.Client, question string) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oLatest,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
			MaxTokens: 4096,
		},
	)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
