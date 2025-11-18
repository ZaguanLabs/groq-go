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

	resp, err := client.Chat.Create(context.Background(), &types.CreateChatCompletionRequest{
		Model: "llama3-8b-8192",
		Messages: []types.ChatCompletionMessageParam{
			{
				Role:    types.RoleUser,
				Content: "Explain quantum computing in one sentence.",
			},
		},
		Temperature: option.Ptr(option.Some(0.5)),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}
