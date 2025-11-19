package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ZaguanLabs/groq-go/groq"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

func main() {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: GROQ_API_KEY environment variable not set")
		os.Exit(1)
	}

	client, err := groq.NewClient(groq.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Testing finish_reason in streaming responses...")
	fmt.Println("================================================")
	fmt.Println()

	// Test with a simple prompt
	stream, err := client.Chat.CreateStream(context.Background(),
		&types.CreateChatCompletionRequest{
			Model: "llama-3.3-70b-versatile",
			Messages: []types.ChatCompletionMessageParam{
				{Role: types.RoleUser, Content: "Say 'Hello World' and nothing else."},
			},
		})
	if err != nil {
		fmt.Printf("Error creating stream: %v\n", err)
		os.Exit(1)
	}
	defer stream.Close()

	sawFinishReason := false
	chunkCount := 0
	var lastChunk *types.ChatCompletionChunk

	fmt.Println("Streaming chunks:")
	for {
		chunk, err := stream.Next(context.Background())
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Printf("\n✓ Received EOF after %d chunks\n", chunkCount)
				break
			}
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		chunkCount++
		lastChunk = chunk

		// Print chunk details
		fmt.Printf("Chunk %d:\n", chunkCount)
		fmt.Printf("  Model: %q\n", chunk.Model)
		fmt.Printf("  Choices: %d\n", len(chunk.Choices))

		for i, choice := range chunk.Choices {
			fmt.Printf("  Choice %d:\n", i)
			fmt.Printf("    Content: %q\n", choice.Delta.Content)
			fmt.Printf("    FinishReason: %q\n", choice.FinishReason)

			if choice.FinishReason != "" {
				sawFinishReason = true
				fmt.Printf("    ✓ FINISH REASON DETECTED: %q\n", choice.FinishReason)
			}
		}

		if chunk.Usage != nil {
			fmt.Printf("  Usage:\n")
			fmt.Printf("    Prompt tokens: %d\n", chunk.Usage.PromptTokens)
			fmt.Printf("    Completion tokens: %d\n", chunk.Usage.CompletionTokens)
			fmt.Printf("    Total tokens: %d\n", chunk.Usage.TotalTokens)
		}
		fmt.Println()
	}

	fmt.Println("\n================================================")
	fmt.Println("Test Results:")
	fmt.Printf("Total chunks received: %d\n", chunkCount)

	if sawFinishReason {
		fmt.Println("✓ SUCCESS: finish_reason was received before EOF")
	} else {
		fmt.Println("✗ FAIL: finish_reason was NOT received before EOF")
		if lastChunk != nil {
			fmt.Println("\nLast chunk details:")
			fmt.Printf("  Model: %q\n", lastChunk.Model)
			fmt.Printf("  Choices: %d\n", len(lastChunk.Choices))
			for i, choice := range lastChunk.Choices {
				fmt.Printf("  Choice %d:\n", i)
				fmt.Printf("    Content: %q\n", choice.Delta.Content)
				fmt.Printf("    FinishReason: %q\n", choice.FinishReason)
			}
		}
		fmt.Println("\nThis indicates the Groq API is not sending finish_reason in streaming mode.")
		fmt.Println("This is a server-side issue, not an SDK bug.")
	}
}
