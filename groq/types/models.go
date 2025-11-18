package types

// Model represents a model
type Model struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

// ModelListResponse represents a list of models
type ModelListResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// ModelDeleted represents a deleted model
type ModelDeleted struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
