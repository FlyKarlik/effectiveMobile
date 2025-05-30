package generics

type ItemsOutput[T any] struct {
	Success bool  `json:"success"`
	Total   int64 `json:"total,omitempty"`
	Items   []T   `json:"items,omitempty"`
	Error   error `json:"error,omitempty"`
}

type ItemOutput[T any] struct {
	Success bool  `json:"success"`
	Item    T     `json:"item,omitempty"`
	Error   error `json:"error,omitempty"`
}
