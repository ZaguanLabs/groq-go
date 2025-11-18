package main

import (
	"context"
	"errors"
	"fmt"
	"io"
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

	stream, err := client.Chat.CreateStream(context.Background(), &types.CreateChatCompletionRequest{
		Model: "llama3-8b-8192",
		Messages: []types.ChatCompletionMessageParam{
			{
				Role:    types.RoleUser,
				Content: "Count from 1 to 5.",
			},
		},
		Temperature: option.Ptr(option.Some(0.5)),
	})
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	fmt.Print("Response: ")
	for {
		chunk, err := stream.Next(context.Background())
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}

		if len(chunk.Choices) > 0 {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
}
