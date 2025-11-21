package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ZaguanLabs/groq-go/groq"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

func main() {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		fmt.Println("GROQ_API_KEY environment variable not set")
		os.Exit(1)
	}

	client, err := groq.NewClient(groq.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		os.Exit(1)
	}

	// Example: Sending a message with document content
	docID := "sales-report-2025"

	req := &types.CreateChatCompletionRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []types.ChatCompletionMessageParam{
			{
				Role: types.RoleUser,
				Content: []interface{}{
					types.ContentPartText{
						Type: "text",
						Text: "Please analyze the following sales data and provide insights:",
					},
					types.ContentPartDocument{
						Type: "document",
						Document: types.ContentPartDocument_Document{
							Data: map[string]interface{}{
								"sales": []interface{}{
									map[string]interface{}{
										"month":   "January",
										"revenue": 125000,
										"units":   450,
									},
									map[string]interface{}{
										"month":   "February",
										"revenue": 142000,
										"units":   520,
									},
									map[string]interface{}{
										"month":   "March",
										"revenue": 138000,
										"units":   495,
									},
								},
								"region":   "North America",
								"currency": "USD",
							},
							ID: &docID,
						},
					},
				},
			},
		},
	}

	fmt.Println("Sending request with document content...")
	fmt.Printf("Document ID: %s\n\n", docID)

	resp, err := client.Chat.Create(context.Background(), req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Response:")
	fmt.Println("─────────────────────────────────────────")
	if len(resp.Choices) > 0 {
		fmt.Println(resp.Choices[0].Message.Content)
	}
	fmt.Println("─────────────────────────────────────────")

	if resp.Usage != nil {
		fmt.Printf("\nUsage:\n")
		fmt.Printf("  Prompt tokens: %d\n", resp.Usage.PromptTokens)
		fmt.Printf("  Completion tokens: %d\n", resp.Usage.CompletionTokens)
		fmt.Printf("  Total tokens: %d\n", resp.Usage.TotalTokens)
	}
}
