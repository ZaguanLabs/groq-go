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

	// Example: Using documents for RAG-like workflows with citations
	documentText := `Quantum computing is a type of computation that harnesses quantum mechanical phenomena.
Unlike classical computers that use bits (0 or 1), quantum computers use quantum bits or qubits.
Qubits can exist in superposition, meaning they can be in multiple states simultaneously.
This property allows quantum computers to process vast amounts of information in parallel.`

	jsonData := map[string]interface{}{
		"title":  "Quantum Computing Basics",
		"author": "Research Team",
		"year":   2024,
	}

	resp, err := client.Chat.Create(context.Background(), &types.CreateChatCompletionRequest{
		Model: string(types.ModelLlama33_70BVersatile),
		Messages: []types.ChatCompletionMessageParam{
			{
				Role:    types.RoleUser,
				Content: "Based on the provided documents, explain what makes quantum computing different from classical computing.",
			},
		},
		Documents: []types.Document{
			{
				ID: func() *string { s := "doc1"; return &s }(),
				Source: &types.DocumentSource{
					Type: "text",
					Text: &documentText,
				},
			},
			{
				ID: func() *string { s := "doc2"; return &s }(),
				Source: &types.DocumentSource{
					Type: "json",
					Data: jsonData,
				},
			},
		},
		CitationOptions: option.Ptr(option.Some("enabled")),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %s\n\n", resp.Choices[0].Message.Content)

	// Display citations
	if len(resp.Choices[0].Message.Annotations) > 0 {
		fmt.Println("Citations:")
		for i, annotation := range resp.Choices[0].Message.Annotations {
			fmt.Printf("  [%d] Type: %s\n", i+1, annotation.Type)
			if annotation.DocumentCitation != nil {
				fmt.Printf("      Document ID: %s\n", annotation.DocumentCitation.DocumentID)
				fmt.Printf("      Text Range: [%d:%d]\n",
					annotation.DocumentCitation.StartIndex,
					annotation.DocumentCitation.EndIndex)
			}
		}
	}
}
