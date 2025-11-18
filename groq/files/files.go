package files

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ZaguanLabs/groq-go/groq/option"
	"github.com/ZaguanLabs/groq-go/groq/types"
)

// Requester defines the interface for sending requests
type Requester interface {
	PostForm(ctx context.Context, path string, formStruct interface{}, result interface{}, opts ...option.RequestOption) error
	Get(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error
	GetStream(ctx context.Context, path string, opts ...option.RequestOption) (*http.Response, error)
	Delete(ctx context.Context, path string, result interface{}, opts ...option.RequestOption) error
}

// Files handles file operations
type Files struct {
	requester Requester
}

// New creates a new Files service
func New(requester Requester) *Files {
	return &Files{requester: requester}
}

// Create uploads a file
func (f *Files) Create(ctx context.Context, req *types.CreateFileRequest, opts ...option.RequestOption) (*types.FileObject, error) {
	var result types.FileObject
	err := f.requester.PostForm(ctx, "/openai/v1/files", req, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// List lists files
func (f *Files) List(ctx context.Context, opts ...option.RequestOption) (*types.FileListResponse, error) {
	var result types.FileListResponse
	err := f.requester.Get(ctx, "/openai/v1/files", &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Retrieve retrieves a file
func (f *Files) Retrieve(ctx context.Context, fileID string, opts ...option.RequestOption) (*types.FileObject, error) {
	var result types.FileObject
	err := f.requester.Get(ctx, fmt.Sprintf("/openai/v1/files/%s", fileID), &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a file
func (f *Files) Delete(ctx context.Context, fileID string, opts ...option.RequestOption) (*types.FileDeleted, error) {
	var result types.FileDeleted
	err := f.requester.Delete(ctx, fmt.Sprintf("/openai/v1/files/%s", fileID), &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Content retrieves the content of a file
func (f *Files) Content(ctx context.Context, fileID string, opts ...option.RequestOption) (io.ReadCloser, error) {
	resp, err := f.requester.GetStream(ctx, fmt.Sprintf("/openai/v1/files/%s/content", fileID), opts...)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
