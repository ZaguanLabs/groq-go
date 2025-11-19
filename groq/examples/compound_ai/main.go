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

	// Example: Using Compound AI with custom model selection and tools
	resp, err := client.Chat.Create(context.Background(), &types.CreateChatCompletionRequest{
		Model: string(types.ModelCompoundBeta),
		Messages: []types.ChatCompletionMessageParam{
			{
				Role:    types.RoleUser,
				Content: "Search for the latest developments in quantum computing and summarize the key findings.",
			},
		},
		CompoundCustom: &types.CompoundCustom{
			Models: &types.CompoundCustomModels{
				AnsweringModel: option.Ptr(option.Some(string(types.ModelLlama33_70BVersatile))),
				ReasoningModel: option.Ptr(option.Some(string(types.ModelQwen332B))),
			},
			Tools: &types.CompoundCustomTools{
				EnabledTools: []string{"web_search"},
			},
		},
		SearchSettings: &types.SearchSettings{
			Country:       option.Ptr(option.Some("united states")),
			IncludeImages: option.Ptr(option.Some(true)),
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %s\n\n", resp.Choices[0].Message.Content)

	// Display executed tools if any
	if len(resp.Choices[0].Message.ExecutedTools) > 0 {
		fmt.Println("Executed Tools:")
		for _, tool := range resp.Choices[0].Message.ExecutedTools {
			fmt.Printf("  - Type: %s\n", tool.Type)
			if tool.SearchResults != nil && len(tool.SearchResults.Results) > 0 {
				fmt.Printf("    Search Results: %d results\n", len(tool.SearchResults.Results))
				for i, result := range tool.SearchResults.Results {
					if i >= 3 {
						break // Show first 3
					}
					if result.Title != nil {
						fmt.Printf("      %d. %s\n", i+1, *result.Title)
					}
					if result.URL != nil {
						fmt.Printf("         URL: %s\n", *result.URL)
					}
				}
			}
		}
	}

	// Display annotations/citations if any
	if len(resp.Choices[0].Message.Annotations) > 0 {
		fmt.Println("\nCitations:")
		for i, annotation := range resp.Choices[0].Message.Annotations {
			fmt.Printf("  [%d] Type: %s\n", i+1, annotation.Type)
			if annotation.FunctionCitation != nil {
				fmt.Printf("      Tool Call ID: %s\n", annotation.FunctionCitation.ToolCallID)
			}
		}
	}
}
