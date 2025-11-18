package types

import (
	"encoding/json"
	"testing"

	"github.com/ZaguanLabs/groq-go/groq/option"
)

func TestChatCompletionRequest_Marshal(t *testing.T) {
	req := CreateChatCompletionRequest{
		Messages: []ChatCompletionMessageParam{
			{Role: RoleUser, Content: "Hello"},
		},
		Model:       "llama3-8b",
		Temperature: option.Ptr(option.Some(0.7)),
	}

	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	expected := `{"messages":[{"role":"user","content":"Hello"}],"model":"llama3-8b","temperature":0.7}`
	if string(b) != expected {
		t.Errorf("Expected %s, got %s", expected, string(b))
	}
}

func TestChatCompletionRequest_Omit(t *testing.T) {
	req := CreateChatCompletionRequest{
		Messages: []ChatCompletionMessageParam{
			{Role: RoleUser, Content: "Hello"},
		},
		Model: "llama3-8b",
		// Temperature unset (nil), should be omitted
	}

	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	expected := `{"messages":[{"role":"user","content":"Hello"}],"model":"llama3-8b"}`
	if string(b) != expected {
		t.Errorf("Expected %s, got %s", expected, string(b))
	}
}

func TestChatCompletionRequest_ExplicitNull(t *testing.T) {
	// Groq API might not support explicit null for temperature, but let's test the mechanism
	req := CreateChatCompletionRequest{
		Messages: []ChatCompletionMessageParam{
			{Role: RoleUser, Content: "Hello"},
		},
		Model:       "llama3-8b",
		Temperature: option.Ptr(option.None[float64]()),
	}

	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	expected := `{"messages":[{"role":"user","content":"Hello"}],"model":"llama3-8b","temperature":null}`
	if string(b) != expected {
		t.Errorf("Expected %s, got %s", expected, string(b))
	}
}

func TestChatCompletion_Unmarshal(t *testing.T) {
	jsonStr := `{
		"id": "chatcmpl-123",
		"object": "chat.completion",
		"created": 1677652288,
		"model": "gpt-3.5-turbo",
		"choices": [{
			"index": 0,
			"message": {
				"role": "assistant",
				"content": "Hello there!"
			},
			"finish_reason": "stop"
		}],
		"usage": {
			"prompt_tokens": 9,
			"completion_tokens": 12,
			"total_tokens": 21
		}
	}`

	var resp ChatCompletion
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if resp.ID != "chatcmpl-123" {
		t.Errorf("ID mismatch")
	}
	if len(resp.Choices) != 1 {
		t.Fatalf("Expected 1 choice")
	}
	if resp.Choices[0].Message.Content != "Hello there!" {
		t.Errorf("Content mismatch")
	}
	if resp.Usage.TotalTokens != 21 {
		t.Errorf("Usage mismatch")
	}
}
