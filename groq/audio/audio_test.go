package audio

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

// mockRequester implements the Requester interface for testing
type mockRequester struct {
	postFunc       func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error
	postStreamFunc func(ctx context.Context, path string, body interface{}, opts ...option.RequestOption) (*http.Response, error)
	postFormFunc   func(ctx context.Context, path string, formStruct interface{}, result interface{}, opts ...option.RequestOption) error
}

func (m *mockRequester) Post(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
	if m.postFunc != nil {
		return m.postFunc(ctx, path, body, result, opts...)
	}
	return nil
}

func (m *mockRequester) PostStream(ctx context.Context, path string, body interface{}, opts ...option.RequestOption) (*http.Response, error) {
	if m.postStreamFunc != nil {
		return m.postStreamFunc(ctx, path, body, opts...)
	}
	return nil, nil
}

func (m *mockRequester) PostForm(ctx context.Context, path string, formStruct interface{}, result interface{}, opts ...option.RequestOption) error {
	if m.postFormFunc != nil {
		return m.postFormFunc(ctx, path, formStruct, result, opts...)
	}
	return nil
}

func TestNew(t *testing.T) {
	mock := &mockRequester{}
	a := New(mock)

	if a == nil {
		t.Fatal("New returned nil")
	}
	if a.Speech == nil {
		t.Error("Speech not initialized")
	}
	if a.Transcriptions == nil {
		t.Error("Transcriptions not initialized")
	}
	if a.Translations == nil {
		t.Error("Translations not initialized")
	}
}

func TestSpeech_Create(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.CreateSpeechRequest
		mockResp    string
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name: "successful speech generation",
			req: &types.CreateSpeechRequest{
				Model: "distil-whisper-large-v3-en",
				Input: "Hello, this is a test.",
				Voice: "alloy",
			},
			mockResp: "fake audio data",
			wantErr:  false,
		},
		{
			name: "with response format",
			req: &types.CreateSpeechRequest{
				Model:          "distil-whisper-large-v3-en",
				Input:          "Test",
				Voice:          "echo",
				ResponseFormat: option.Ptr(option.Some("mp3")),
			},
			mockResp: "mp3 audio data",
			wantErr:  false,
		},
		{
			name: "with speed",
			req: &types.CreateSpeechRequest{
				Model: "distil-whisper-large-v3-en",
				Input: "Fast speech",
				Voice: "fable",
				Speed: option.Ptr(option.Some(1.5)),
			},
			mockResp: "fast audio data",
			wantErr:  false,
		},
		{
			name: "error from requester",
			req: &types.CreateSpeechRequest{
				Model: "distil-whisper-large-v3-en",
				Input: "Test",
				Voice: "alloy",
			},
			mockErr:     errors.New("generation error"),
			wantErr:     true,
			errContains: "generation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRequester{
				postStreamFunc: func(ctx context.Context, path string, body interface{}, opts ...option.RequestOption) (*http.Response, error) {
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					// Verify path
					if path != "/openai/v1/audio/speech" {
						t.Errorf("unexpected path: %s", path)
					}
					// Return mock audio data
					resp := &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(tt.mockResp)),
						Header:     make(http.Header),
					}
					resp.Header.Set("Content-Type", "audio/mpeg")
					return resp, nil
				},
			}

			a := New(mock)
			reader, err := a.Speech.Create(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if reader == nil {
				t.Fatal("reader is nil")
			}

			// Read and verify data
			data, err := io.ReadAll(reader)
			if err != nil {
				t.Fatalf("failed to read audio data: %v", err)
			}

			if string(data) != tt.mockResp {
				t.Errorf("audio data = %q, want %q", string(data), tt.mockResp)
			}

			reader.Close()
		})
	}
}

func TestTranscriptions_Create(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.CreateTranscriptionRequest
		mockResp    *types.Transcription
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name: "successful transcription",
			req: &types.CreateTranscriptionRequest{
				File:  strings.NewReader("fake audio data"),
				Model: "whisper-large-v3",
			},
			mockResp: &types.Transcription{
				Text: "This is the transcribed text.",
			},
			wantErr: false,
		},
		{
			name: "with language",
			req: &types.CreateTranscriptionRequest{
				File:     strings.NewReader("audio data"),
				Model:    "whisper-large-v3",
				Language: option.Ptr(option.Some("en")),
			},
			mockResp: &types.Transcription{
				Text: "English transcription.",
			},
			wantErr: false,
		},
		{
			name: "with prompt",
			req: &types.CreateTranscriptionRequest{
				File:   strings.NewReader("audio data"),
				Model:  "whisper-large-v3",
				Prompt: option.Ptr(option.Some("Context prompt")),
			},
			mockResp: &types.Transcription{
				Text: "Transcription with context.",
			},
			wantErr: false,
		},
		{
			name: "with temperature",
			req: &types.CreateTranscriptionRequest{
				File:        strings.NewReader("audio data"),
				Model:       "whisper-large-v3",
				Temperature: option.Ptr(option.Some(0.5)),
			},
			mockResp: &types.Transcription{
				Text: "Temperature-controlled transcription.",
			},
			wantErr: false,
		},
		{
			name: "error from requester",
			req: &types.CreateTranscriptionRequest{
				File:  strings.NewReader("bad audio"),
				Model: "whisper-large-v3",
			},
			mockErr:     errors.New("transcription failed"),
			wantErr:     true,
			errContains: "transcription failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRequester{
				postFormFunc: func(ctx context.Context, path string, formStruct interface{}, result interface{}, opts ...option.RequestOption) error {
					if tt.mockErr != nil {
						return tt.mockErr
					}
					if tt.mockResp != nil {
						respBytes, _ := json.Marshal(tt.mockResp)
						json.Unmarshal(respBytes, result)
					}
					// Verify path
					if path != "/openai/v1/audio/transcriptions" {
						t.Errorf("unexpected path: %s", path)
					}
					return nil
				},
			}

			a := New(mock)
			resp, err := a.Transcriptions.Create(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp == nil {
				t.Fatal("response is nil")
			}

			if resp.Text != tt.mockResp.Text {
				t.Errorf("text = %q, want %q", resp.Text, tt.mockResp.Text)
			}
		})
	}
}

func TestTranslations_Create(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.CreateTranslationRequest
		mockResp    *types.Translation
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name: "successful translation",
			req: &types.CreateTranslationRequest{
				File:  strings.NewReader("fake audio data"),
				Model: "whisper-large-v3",
			},
			mockResp: &types.Translation{
				Text: "This is the translated text in English.",
			},
			wantErr: false,
		},
		{
			name: "with prompt",
			req: &types.CreateTranslationRequest{
				File:   strings.NewReader("audio data"),
				Model:  "whisper-large-v3",
				Prompt: option.Ptr(option.Some("Translation context")),
			},
			mockResp: &types.Translation{
				Text: "Translated with context.",
			},
			wantErr: false,
		},
		{
			name: "with response format",
			req: &types.CreateTranslationRequest{
				File:           strings.NewReader("audio data"),
				Model:          "whisper-large-v3",
				ResponseFormat: option.Ptr(option.Some("text")),
			},
			mockResp: &types.Translation{
				Text: "Plain text translation.",
			},
			wantErr: false,
		},
		{
			name: "error from requester",
			req: &types.CreateTranslationRequest{
				File:  strings.NewReader("bad audio"),
				Model: "whisper-large-v3",
			},
			mockErr:     errors.New("translation failed"),
			wantErr:     true,
			errContains: "translation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRequester{
				postFormFunc: func(ctx context.Context, path string, formStruct interface{}, result interface{}, opts ...option.RequestOption) error {
					if tt.mockErr != nil {
						return tt.mockErr
					}
					if tt.mockResp != nil {
						respBytes, _ := json.Marshal(tt.mockResp)
						json.Unmarshal(respBytes, result)
					}
					// Verify path
					if path != "/openai/v1/audio/translations" {
						t.Errorf("unexpected path: %s", path)
					}
					return nil
				},
			}

			a := New(mock)
			resp, err := a.Translations.Create(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp == nil {
				t.Fatal("response is nil")
			}

			if resp.Text != tt.mockResp.Text {
				t.Errorf("text = %q, want %q", resp.Text, tt.mockResp.Text)
			}
		})
	}
}

func TestAudio_WithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	t.Run("Speech with cancelled context", func(t *testing.T) {
		mock := &mockRequester{
			postStreamFunc: func(ctx context.Context, path string, body interface{}, opts ...option.RequestOption) (*http.Response, error) {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
					return nil, nil
				}
			},
		}

		a := New(mock)
		req := &types.CreateSpeechRequest{
			Model: "distil-whisper-large-v3-en",
			Input: "Test",
			Voice: "alloy",
		}

		_, err := a.Speech.Create(ctx, req)
		if err == nil {
			t.Fatal("expected context cancellation error")
		}
	})

	t.Run("Transcriptions with cancelled context", func(t *testing.T) {
		mock := &mockRequester{
			postFormFunc: func(ctx context.Context, path string, formStruct interface{}, result interface{}, opts ...option.RequestOption) error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return nil
				}
			},
		}

		a := New(mock)
		req := &types.CreateTranscriptionRequest{
			File:  strings.NewReader("audio"),
			Model: "whisper-large-v3",
		}

		_, err := a.Transcriptions.Create(ctx, req)
		if err == nil {
			t.Fatal("expected context cancellation error")
		}
	})

	t.Run("Translations with cancelled context", func(t *testing.T) {
		mock := &mockRequester{
			postFormFunc: func(ctx context.Context, path string, formStruct interface{}, result interface{}, opts ...option.RequestOption) error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return nil
				}
			},
		}

		a := New(mock)
		req := &types.CreateTranslationRequest{
			File:  strings.NewReader("audio"),
			Model: "whisper-large-v3",
		}

		_, err := a.Translations.Create(ctx, req)
		if err == nil {
			t.Fatal("expected context cancellation error")
		}
	})
}
