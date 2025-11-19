package files

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

type mockRequester struct {
	postFormFunc  func(ctx context.Context, path string, formStruct interface{}, result interface{}, opts ...option.RequestOption) error
	getFunc       func(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error
	getStreamFunc func(ctx context.Context, path string, opts ...option.RequestOption) (*http.Response, error)
	deleteFunc    func(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error
}

func (m *mockRequester) PostForm(ctx context.Context, path string, formStruct interface{}, result interface{}, opts ...option.RequestOption) error {
	if m.postFormFunc != nil {
		return m.postFormFunc(ctx, path, formStruct, result, opts...)
	}
	return nil
}

func (m *mockRequester) Get(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
	if m.getFunc != nil {
		return m.getFunc(ctx, path, result, opts...)
	}
	return nil
}

func (m *mockRequester) GetStream(ctx context.Context, path string, opts ...option.RequestOption) (*http.Response, error) {
	if m.getStreamFunc != nil {
		return m.getStreamFunc(ctx, path, opts...)
	}
	return nil, nil
}

func (m *mockRequester) Delete(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, path, result, opts...)
	}
	return nil
}

func TestNew(t *testing.T) {
	mock := &mockRequester{}
	f := New(mock)

	if f == nil {
		t.Fatal("New returned nil")
	}
	if f.requester != mock {
		t.Error("requester not set correctly")
	}
}

func TestFiles_Create(t *testing.T) {
	tests := []struct {
		name        string
		req         *types.CreateFileRequest
		mockResp    *types.FileObject
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name: "successful upload",
			req: &types.CreateFileRequest{
				File:    strings.NewReader("file content"),
				Purpose: "batch",
			},
			mockResp: &types.FileObject{
				ID:        "file-123",
				Bytes:     100,
				CreatedAt: 1234567890,
				Filename:  "test.jsonl",
				Object:    "file",
				Purpose:   "batch",
			},
			wantErr: false,
		},
		{
			name: "error from requester",
			req: &types.CreateFileRequest{
				File:    strings.NewReader("content"),
				Purpose: "batch",
			},
			mockErr:     errors.New("upload failed"),
			wantErr:     true,
			errContains: "upload failed",
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
					if path != "/openai/v1/files" {
						t.Errorf("unexpected path: %s", path)
					}
					return nil
				},
			}

			f := New(mock)
			resp, err := f.Create(context.Background(), tt.req)

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

func TestFiles_List(t *testing.T) {
	tests := []struct {
		name        string
		mockResp    *types.FileListResponse
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name: "successful list",
			mockResp: &types.FileListResponse{
				Object: "list",
				Data: []types.FileObject{
					{ID: "file-1", Filename: "file1.jsonl"},
					{ID: "file-2", Filename: "file2.jsonl"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty list",
			mockResp: &types.FileListResponse{
				Object: "list",
				Data:   []types.FileObject{},
			},
			wantErr: false,
		},
		{
			name:        "error from requester",
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
					if path != "/openai/v1/files" {
						t.Errorf("unexpected path: %s", path)
					}
					return nil
				},
			}

			f := New(mock)
			resp, err := f.List(context.Background())

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
				t.Errorf("got %d files, want %d", len(resp.Data), len(tt.mockResp.Data))
			}
		})
	}
}

func TestFiles_Retrieve(t *testing.T) {
	tests := []struct {
		name        string
		fileID      string
		mockResp    *types.FileObject
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name:   "successful retrieve",
			fileID: "file-123",
			mockResp: &types.FileObject{
				ID:       "file-123",
				Filename: "test.jsonl",
			},
			wantErr: false,
		},
		{
			name:        "error from requester",
			fileID:      "invalid",
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
					expectedPath := "/openai/v1/files/" + tt.fileID
					if path != expectedPath {
						t.Errorf("path = %q, want %q", path, expectedPath)
					}
					return nil
				},
			}

			f := New(mock)
			resp, err := f.Retrieve(context.Background(), tt.fileID)

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

func TestFiles_Delete(t *testing.T) {
	tests := []struct {
		name        string
		fileID      string
		mockResp    *types.FileDeleted
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name:   "successful delete",
			fileID: "file-123",
			mockResp: &types.FileDeleted{
				ID:      "file-123",
				Object:  "file",
				Deleted: true,
			},
			wantErr: false,
		},
		{
			name:        "error from requester",
			fileID:      "invalid",
			mockErr:     errors.New("delete failed"),
			wantErr:     true,
			errContains: "delete failed",
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
					expectedPath := "/openai/v1/files/" + tt.fileID
					if path != expectedPath {
						t.Errorf("path = %q, want %q", path, expectedPath)
					}
					return nil
				},
			}

			f := New(mock)
			resp, err := f.Delete(context.Background(), tt.fileID)

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

			if !resp.Deleted {
				t.Error("file not marked as deleted")
			}
		})
	}
}

func TestFiles_Content(t *testing.T) {
	tests := []struct {
		name        string
		fileID      string
		mockContent string
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name:        "successful content retrieval",
			fileID:      "file-123",
			mockContent: "file content here",
			wantErr:     false,
		},
		{
			name:        "error from requester",
			fileID:      "invalid",
			mockErr:     errors.New("content not found"),
			wantErr:     true,
			errContains: "content not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRequester{
				getStreamFunc: func(ctx context.Context, path string, opts ...option.RequestOption) (*http.Response, error) {
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					expectedPath := "/openai/v1/files/" + tt.fileID + "/content"
					if path != expectedPath {
						t.Errorf("path = %q, want %q", path, expectedPath)
					}
					resp := &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(tt.mockContent)),
						Header:     make(http.Header),
					}
					return resp, nil
				},
			}

			f := New(mock)
			reader, err := f.Content(context.Background(), tt.fileID)

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

			content, err := io.ReadAll(reader)
			if err != nil {
				t.Fatalf("failed to read content: %v", err)
			}

			if string(content) != tt.mockContent {
				t.Errorf("content = %q, want %q", string(content), tt.mockContent)
			}

			reader.Close()
		})
	}
}

func TestFiles_WithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

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

	f := New(mock)
	req := &types.CreateFileRequest{
		File:    strings.NewReader("test"),
		Purpose: "batch",
	}

	_, err := f.Create(ctx, req)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}
