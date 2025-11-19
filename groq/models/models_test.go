package models

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

// mockRequester implements the Requester interface for testing
type mockRequester struct {
	getFunc    func(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error
	deleteFunc func(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error
}

func (m *mockRequester) Get(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
	if m.getFunc != nil {
		return m.getFunc(ctx, path, result, opts...)
	}
	return nil
}

func (m *mockRequester) Delete(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, path, result, opts...)
	}
	return nil
}

func TestNew(t *testing.T) {
	mock := &mockRequester{}
	m := New(mock)

	if m == nil {
		t.Fatal("New returned nil")
	}
	if m.requester != mock {
		t.Error("requester not set correctly")
	}
}

func TestModels_List(t *testing.T) {
	tests := []struct {
		name        string
		mockResp    *types.ModelListResponse
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name: "successful list",
			mockResp: &types.ModelListResponse{
				Object: "list",
				Data: []types.Model{
					{
						ID:      "llama-3.1-8b-instant",
						Created: 1234567890,
						Object:  "model",
						OwnedBy: "groq",
					},
					{
						ID:      "llama-3.3-70b-versatile",
						Created: 1234567891,
						Object:  "model",
						OwnedBy: "groq",
					},
				},
			},
			wantErr: false,
		},
		{
			name:        "error from requester",
			mockErr:     errors.New("network error"),
			wantErr:     true,
			errContains: "network error",
		},
		{
			name: "empty list",
			mockResp: &types.ModelListResponse{
				Object: "list",
				Data:   []types.Model{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRequester{
				getFunc: func(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
					if tt.mockErr != nil {
						return tt.mockErr
					}
					if tt.mockResp != nil {
						respBytes, _ := json.Marshal(tt.mockResp)
						json.Unmarshal(respBytes, result)
					}
					// Verify path
					if path != "/openai/v1/models" {
						t.Errorf("unexpected path: %s", path)
					}
					return nil
				},
			}

			m := New(mock)
			resp, err := m.List(context.Background())

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

			if len(resp.Data) != len(tt.mockResp.Data) {
				t.Errorf("got %d models, want %d", len(resp.Data), len(tt.mockResp.Data))
			}
		})
	}
}

func TestModels_Retrieve(t *testing.T) {
	tests := []struct {
		name        string
		modelID     string
		mockResp    *types.Model
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name:    "successful retrieve",
			modelID: "llama-3.1-8b-instant",
			mockResp: &types.Model{
				ID:      "llama-3.1-8b-instant",
				Created: 1234567890,
				Object:  "model",
				OwnedBy: "groq",
			},
			wantErr: false,
		},
		{
			name:        "error from requester",
			modelID:     "invalid-model",
			mockErr:     errors.New("model not found"),
			wantErr:     true,
			errContains: "model not found",
		},
		{
			name:    "different model",
			modelID: "llama-3.3-70b-versatile",
			mockResp: &types.Model{
				ID:      "llama-3.3-70b-versatile",
				Created: 1234567891,
				Object:  "model",
				OwnedBy: "groq",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRequester{
				getFunc: func(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
					if tt.mockErr != nil {
						return tt.mockErr
					}
					if tt.mockResp != nil {
						respBytes, _ := json.Marshal(tt.mockResp)
						json.Unmarshal(respBytes, result)
					}
					// Verify path
					expectedPath := "/openai/v1/models/" + tt.modelID
					if path != expectedPath {
						t.Errorf("path = %q, want %q", path, expectedPath)
					}
					return nil
				},
			}

			m := New(mock)
			resp, err := m.Retrieve(context.Background(), tt.modelID)

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

			if resp.ID != tt.mockResp.ID {
				t.Errorf("ID = %q, want %q", resp.ID, tt.mockResp.ID)
			}
		})
	}
}

func TestModels_Delete(t *testing.T) {
	tests := []struct {
		name        string
		modelID     string
		mockResp    *types.ModelDeleted
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name:    "successful delete",
			modelID: "ft-model-123",
			mockResp: &types.ModelDeleted{
				ID:      "ft-model-123",
				Object:  "model",
				Deleted: true,
			},
			wantErr: false,
		},
		{
			name:        "error from requester",
			modelID:     "invalid-model",
			mockErr:     errors.New("cannot delete model"),
			wantErr:     true,
			errContains: "cannot delete model",
		},
		{
			name:    "delete not allowed",
			modelID: "base-model",
			mockResp: &types.ModelDeleted{
				ID:      "base-model",
				Object:  "model",
				Deleted: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRequester{
				deleteFunc: func(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
					if tt.mockErr != nil {
						return tt.mockErr
					}
					if tt.mockResp != nil {
						respBytes, _ := json.Marshal(tt.mockResp)
						json.Unmarshal(respBytes, result)
					}
					// Verify path
					expectedPath := "/openai/v1/models/" + tt.modelID
					if path != expectedPath {
						t.Errorf("path = %q, want %q", path, expectedPath)
					}
					return nil
				},
			}

			m := New(mock)
			resp, err := m.Delete(context.Background(), tt.modelID)

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

			if resp.ID != tt.mockResp.ID {
				t.Errorf("ID = %q, want %q", resp.ID, tt.mockResp.ID)
			}

			if resp.Deleted != tt.mockResp.Deleted {
				t.Errorf("Deleted = %v, want %v", resp.Deleted, tt.mockResp.Deleted)
			}
		})
	}
}

func TestModels_ListWithOptions(t *testing.T) {
	mock := &mockRequester{
		getFunc: func(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
			if len(opts) == 0 {
				t.Error("expected options to be passed")
			}
			return nil
		},
	}

	m := New(mock)
	_, err := m.List(
		context.Background(),
		option.WithRequestHeader("X-Custom", "test"),
		option.WithRequestQuery("limit", "10"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestModels_RetrieveWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	mock := &mockRequester{
		getFunc: func(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		},
	}

	m := New(mock)
	_, err := m.Retrieve(ctx, "test-model")
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestModels_DeleteWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	mock := &mockRequester{
		deleteFunc: func(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		},
	}

	m := New(mock)
	_, err := m.Delete(ctx, "test-model")
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}
