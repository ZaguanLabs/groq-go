package types

import (
	"encoding/json"
	"testing"
)

func TestContentPartText_JSON(t *testing.T) {
	part := ContentPartText{
		Type: "text",
		Text: "Hello, world!",
	}

	data, err := json.Marshal(part)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := `{"type":"text","text":"Hello, world!"}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", string(data), expected)
	}

	var decoded ContentPartText
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Type != part.Type || decoded.Text != part.Text {
		t.Errorf("decoded mismatch: got %+v, want %+v", decoded, part)
	}
}

func TestContentPartImage_JSON(t *testing.T) {
	part := ContentPartImage{
		Type: "image_url",
		ImageURL: ContentPartImage_ImageURL{
			URL:    "https://example.com/image.jpg",
			Detail: "high",
		},
	}

	data, err := json.Marshal(part)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ContentPartImage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Type != part.Type {
		t.Errorf("type mismatch: got %s, want %s", decoded.Type, part.Type)
	}
	if decoded.ImageURL.URL != part.ImageURL.URL {
		t.Errorf("URL mismatch: got %s, want %s", decoded.ImageURL.URL, part.ImageURL.URL)
	}
	if decoded.ImageURL.Detail != part.ImageURL.Detail {
		t.Errorf("detail mismatch: got %s, want %s", decoded.ImageURL.Detail, part.ImageURL.Detail)
	}
}

func TestContentPartDocument_JSON(t *testing.T) {
	id := "doc-123"
	part := ContentPartDocument{
		Type: "document",
		Document: ContentPartDocument_Document{
			Data: map[string]interface{}{
				"sales":  []interface{}{100.0, 200.0, 300.0},
				"region": "North America",
				"metadata": map[string]interface{}{
					"year":    2025.0,
					"quarter": "Q4",
				},
			},
			ID: &id,
		},
	}

	data, err := json.Marshal(part)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ContentPartDocument
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Type != part.Type {
		t.Errorf("type mismatch: got %s, want %s", decoded.Type, part.Type)
	}

	if decoded.Document.ID == nil || *decoded.Document.ID != *part.Document.ID {
		t.Errorf("ID mismatch: got %v, want %v", decoded.Document.ID, part.Document.ID)
	}

	// Verify data structure
	if decoded.Document.Data["region"] != "North America" {
		t.Errorf("region mismatch: got %v, want North America", decoded.Document.Data["region"])
	}

	sales, ok := decoded.Document.Data["sales"].([]interface{})
	if !ok {
		t.Fatalf("sales is not a slice")
	}
	if len(sales) != 3 {
		t.Errorf("sales length: got %d, want 3", len(sales))
	}
}

func TestContentPartDocument_WithoutID(t *testing.T) {
	part := ContentPartDocument{
		Type: "document",
		Document: ContentPartDocument_Document{
			Data: map[string]interface{}{
				"key": "value",
			},
			ID: nil,
		},
	}

	data, err := json.Marshal(part)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify ID is omitted when nil
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	doc, ok := raw["document"].(map[string]interface{})
	if !ok {
		t.Fatalf("document is not a map")
	}

	if _, hasID := doc["id"]; hasID {
		t.Errorf("id should be omitted when nil, but was present")
	}
}

func TestMultimodalMessage_JSON(t *testing.T) {
	id := "doc-456"
	msg := ChatCompletionMessageParam{
		Role: RoleUser,
		Content: []interface{}{
			ContentPartText{
				Type: "text",
				Text: "Analyze this data:",
			},
			ContentPartDocument{
				Type: "document",
				Document: ContentPartDocument_Document{
					Data: map[string]interface{}{
						"values": []interface{}{1.0, 2.0, 3.0},
					},
					ID: &id,
				},
			},
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ChatCompletionMessageParam
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Role != msg.Role {
		t.Errorf("role mismatch: got %s, want %s", decoded.Role, msg.Role)
	}

	// Content should be preserved as slice
	content, ok := decoded.Content.([]interface{})
	if !ok {
		t.Fatalf("content is not a slice")
	}

	if len(content) != 2 {
		t.Errorf("content length: got %d, want 2", len(content))
	}
}
