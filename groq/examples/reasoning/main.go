package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ZaguanLabs/groq-go/groq"
	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

func main() {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set GROQ_API_KEY environment variable")
		return
	}

	client, err := groq.NewClient(groq.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}

	// Example: Using reasoning models with parsed reasoning output
	resp, err := client.Chat.Create(context.Background(), &types.CreateChatCompletionRequest{
		Model: string(types.ModelQwen332B),
		Messages: []types.ChatCompletionMessageParam{
			{
				Role:    types.RoleUser,
				Content: "Solve this logic puzzle: If all roses are flowers, and some flowers fade quickly, can we conclude that some roses fade quickly?",
			},
		},
		ReasoningEffort:  option.Ptr(option.Some("medium")),
		ReasoningFormat:  option.Ptr(option.Some("parsed")),
		IncludeReasoning: option.Ptr(option.Some(true)),
		Temperature:      option.Ptr(option.Some(0.3)),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %s\n\n", resp.Choices[0].Message.Content)

	// Display reasoning if available
	if resp.Choices[0].Message.Reasoning != nil {
		fmt.Println("Model's Reasoning Process:")
		fmt.Println(*resp.Choices[0].Message.Reasoning)
	}
}
