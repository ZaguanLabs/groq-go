package types

// ContentPart represents a part of message content (text, image, or document)
type ContentPart interface {
	contentPart()
}

// ContentPartText represents a text content part
type ContentPartText struct {
	Type string `json:"type"` // "text"
	Text string `json:"text"`
}

func (ContentPartText) contentPart() {}

// ContentPartImage represents an image content part
type ContentPartImage struct {
	Type     string           `json:"type"` // "image_url"
	ImageURL ContentPartImage_ImageURL `json:"image_url"`
}

func (ContentPartImage) contentPart() {}

// ContentPartImage_ImageURL represents the image URL details
type ContentPartImage_ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"` // "auto", "low", "high"
}

// ContentPartDocument represents a document content part
type ContentPartDocument struct {
	Type     string                       `json:"type"` // "document"
	Document ContentPartDocument_Document `json:"document"`
}

func (ContentPartDocument) contentPart() {}

// ContentPartDocument_Document represents the document details
type ContentPartDocument_Document struct {
	Data map[string]interface{} `json:"data"` // The JSON document data
	ID   *string                `json:"id,omitempty"` // Optional unique identifier for the document
}
