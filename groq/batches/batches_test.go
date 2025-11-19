package batches

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

type mockRequester struct {
	postFunc func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error
	getFunc  func(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error
}

func (m *mockRequester) Post(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
	if m.postFunc != nil {
		return m.postFunc(ctx, path, body, result, opts...)
	}
	return nil
}

func (m *mockRequester) Get(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
	if m.getFunc != nil {
		return m.getFunc(ctx, path, result, opts...)
	}
	return nil
}

func TestNew(t *testing.T) {
	mock := &mockRequester{}
	b := New(mock)

	if b == nil {
		t.Fatal("New returned nil")
	}
	if b.requester != mock {
		t.Error("requester not set correctly")
	}
}

func TestBatches_Create(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.CreateBatchRequest
		mockResp    *types.Batch
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name: "successful create",
			req: &types.CreateBatchRequest{
				InputFileID:      "file-123",
				Endpoint:         "/v1/chat/completions",
				CompletionWindow: "24h",
			},
			mockResp: &types.Batch{
				ID:               "batch-123",
				Object:           "batch",
				Endpoint:         "/v1/chat/completions",
				InputFileID:      "file-123",
				CompletionWindow: "24h",
				Status:           "validating",
				CreatedAt:        1234567890,
			},
			wantErr: false,
		},
		{
			name: "with metadata",
			req: &types.CreateBatchRequest{
				InputFileID:      "file-456",
				Endpoint:         "/v1/chat/completions",
				CompletionWindow: "24h",
				Metadata:         map[string]string{"key": "value"},
			},
			mockResp: &types.Batch{
				ID:       "batch-456",
				Status:   "validating",
				Metadata: map[string]string{"key": "value"},
			},
			wantErr: false,
		},
		{
			name: "error from requester",
			req: &types.CreateBatchRequest{
				InputFileID:      "file-bad",
				Endpoint:         "/v1/chat/completions",
				CompletionWindow: "24h",
			},
			mockErr:     errors.New("create failed"),
			wantErr:     true,
			errContains: "create failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRequester{
				postFunc: func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
					if tt.mockErr != nil {
						return tt.mockErr
					}
					if tt.mockResp != nil {
						respBytes, _ := json.Marshal(tt.mockResp)
						json.Unmarshal(respBytes, result)
					}
					if path != "/openai/v1/batches" {
						t.Errorf("unexpected path: %s", path)
					}
					return nil
				},
			}

			b := New(mock)
			resp, err := b.Create(context.Background(), tt.req)

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

func TestBatches_Retrieve(t *testing.T) {
	tests := []struct {
		name        string
		batchID     string
		mockResp    *types.Batch
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name:    "successful retrieve",
			batchID: "batch-123",
			mockResp: &types.Batch{
				ID:     "batch-123",
				Status: "completed",
			},
			wantErr: false,
		},
		{
			name:        "error from requester",
			batchID:     "invalid",
			mockErr:     errors.New("not found"),
			wantErr:     true,
			errContains: "not found",
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
					expectedPath := "/openai/v1/batches/" + tt.batchID
					if path != expectedPath {
						t.Errorf("path = %q, want %q", path, expectedPath)
					}
					return nil
				},
			}

			b := New(mock)
			resp, err := b.Retrieve(context.Background(), tt.batchID)

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

func TestBatches_Cancel(t *testing.T) {
	tests := []struct {
		name        string
		batchID     string
		mockResp    *types.Batch
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name:    "successful cancel",
			batchID: "batch-123",
			mockResp: &types.Batch{
				ID:     "batch-123",
				Status: "cancelling",
			},
			wantErr: false,
		},
		{
			name:        "error from requester",
			batchID:     "invalid",
			mockErr:     errors.New("cancel failed"),
			wantErr:     true,
			errContains: "cancel failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRequester{
				postFunc: func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
					if tt.mockErr != nil {
						return tt.mockErr
					}
					if tt.mockResp != nil {
						respBytes, _ := json.Marshal(tt.mockResp)
						json.Unmarshal(respBytes, result)
					}
					expectedPath := "/openai/v1/batches/" + tt.batchID + "/cancel"
					if path != expectedPath {
						t.Errorf("path = %q, want %q", path, expectedPath)
					}
					return nil
				},
			}

			b := New(mock)
			resp, err := b.Cancel(context.Background(), tt.batchID)

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

			if resp.Status != "cancelling" {
				t.Errorf("status = %q, want %q", resp.Status, "cancelling")
			}
		})
	}
}

func TestBatches_List(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.ListBatchesRequest
		mockResp    *types.BatchListResponse
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name: "list without params",
			req:  nil,
			mockResp: &types.BatchListResponse{
				Object: "list",
				Data: []types.Batch{
					{ID: "batch-1", Status: "completed"},
					{ID: "batch-2", Status: "in_progress"},
				},
			},
			wantErr: false,
		},
		{
			name: "list with limit",
			req: &types.ListBatchesRequest{
				Limit: option.Ptr(option.Some(10)),
			},
			mockResp: &types.BatchListResponse{
				Object: "list",
				Data:   []types.Batch{{ID: "batch-1"}},
			},
			wantErr: false,
		},
		{
			name: "list with after",
			req: &types.ListBatchesRequest{
				After: option.Ptr(option.Some("batch-100")),
			},
			mockResp: &types.BatchListResponse{
				Object: "list",
				Data:   []types.Batch{{ID: "batch-101"}},
			},
			wantErr: false,
		},
		{
			name:        "error from requester",
			req:         nil,
			mockErr:     errors.New("list failed"),
			wantErr:     true,
			errContains: "list failed",
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
					if path != "/openai/v1/batches" {
						t.Errorf("unexpected path: %s", path)
					}
					return nil
				},
			}

			b := New(mock)
			resp, err := b.List(context.Background(), tt.req)

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
				t.Errorf("got %d batches, want %d", len(resp.Data), len(tt.mockResp.Data))
			}
		})
	}
}

func TestBatches_WithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mock := &mockRequester{
		postFunc: func(ctx context.Context, path string, body, result interface{}, opts ...option.RequestOption) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		},
	}

	b := New(mock)
	req := &types.CreateBatchRequest{
		InputFileID:      "file-123",
		Endpoint:         "/v1/chat/completions",
		CompletionWindow: "24h",
	}

	_, err := b.Create(ctx, req)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}
