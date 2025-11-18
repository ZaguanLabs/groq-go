package types

import (
	"github.com/ZaguanLabs/groq-go/groq/option"
)

// Transcription represents a transcription response
type Transcription struct {
	Text string `json:"text"`
}

// Translation represents a translation response
type Translation struct {
	Text string `json:"text"`
}

// CreateSpeechRequest represents request parameters for speech generation
type CreateSpeechRequest struct {
	Model          string                    `json:"model"`
	Input          string                    `json:"input"`
	Voice          string                    `json:"voice"`
	ResponseFormat *option.Optional[string]  `json:"response_format,omitempty"` // mp3, opus, aac, flac
	Speed          *option.Optional[float64] `json:"speed,omitempty"`
}

// CreateTranscriptionRequest represents request parameters for transcription
// Note: This struct is for parameters, but actual request is multipart/form-data.
// The client will need to handle file fields separately or via this struct with special handling.
type CreateTranscriptionRequest struct {
	File                   interface{}               `json:"file"` // io.Reader or filename string
	Model                  string                    `json:"model"`
	Language               *option.Optional[string]  `json:"language,omitempty"`
	Prompt                 *option.Optional[string]  `json:"prompt,omitempty"`
	ResponseFormat         *option.Optional[string]  `json:"response_format,omitempty"` // json, text, srt, verbose_json, vtt
	Temperature            *option.Optional[float64] `json:"temperature,omitempty"`
	TimestampGranularities []string                  `json:"timestamp_granularities[],omitempty"` // array in query/form?
}

// CreateTranslationRequest represents request parameters for translation
type CreateTranslationRequest struct {
	File           interface{}               `json:"file"`
	Model          string                    `json:"model"`
	Prompt         *option.Optional[string]  `json:"prompt,omitempty"`
	ResponseFormat *option.Optional[string]  `json:"response_format,omitempty"`
	Temperature    *option.Optional[float64] `json:"temperature,omitempty"`
}
