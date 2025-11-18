package audio

import (
	"context"
	"io"

	"net/http"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

// Requester defines the interface for sending requests
type Requester interface {
	Post(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error
	PostStream(ctx context.Context, path string, body interface{}, opts ...option.RequestOption) (*http.Response, error)
	PostForm(ctx context.Context, path string, formStruct interface{}, result interface{}, opts ...option.RequestOption) error
}

// Audio handles audio requests
type Audio struct {
	Speech         *Speech
	Transcriptions *Transcriptions
	Translations   *Translations
}

// New creates a new Audio service
func New(requester Requester) *Audio {
	return &Audio{
		Speech:         &Speech{requester: requester},
		Transcriptions: &Transcriptions{requester: requester},
		Translations:   &Translations{requester: requester},
	}
}

type Speech struct {
	requester Requester
}

// Create generates audio from the input text
func (s *Speech) Create(ctx context.Context, req *types.CreateSpeechRequest, opts ...option.RequestOption) (io.ReadCloser, error) {
	// Speech returns binary data, so we use PostStream to get the raw response
	resp, err := s.requester.PostStream(ctx, "/openai/v1/audio/speech", req, opts...)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

type Transcriptions struct {
	requester Requester
}

// Create transcribes audio into the input language
func (t *Transcriptions) Create(ctx context.Context, req *types.CreateTranscriptionRequest, opts ...option.RequestOption) (*types.Transcription, error) {
	var result types.Transcription
	err := t.requester.PostForm(ctx, "/openai/v1/audio/transcriptions", req, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type Translations struct {
	requester Requester
}

// Create translates audio into English
func (t *Translations) Create(ctx context.Context, req *types.CreateTranslationRequest, opts ...option.RequestOption) (*types.Translation, error) {
	var result types.Translation
	err := t.requester.PostForm(ctx, "/openai/v1/audio/translations", req, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
