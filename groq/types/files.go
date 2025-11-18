package types

// FileObject represents a file
type FileObject struct {
	ID        string `json:"id"`
	Bytes     int64  `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Object    string `json:"object"`
	Purpose   string `json:"purpose"`
}

// FileListResponse represents a list of files
type FileListResponse struct {
	Object string       `json:"object"`
	Data   []FileObject `json:"data"`
}

// FileDeleted represents a deleted file response
type FileDeleted struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// CreateFileRequest represents parameters to upload a file
type CreateFileRequest struct {
	File    interface{} `json:"file"` // io.Reader or file path
	Purpose string      `json:"purpose"`
}
