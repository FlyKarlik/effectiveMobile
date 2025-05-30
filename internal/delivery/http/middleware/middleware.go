package http_middleware

type HTTPMiddleware struct{}

func New() *HTTPMiddleware {
	return &HTTPMiddleware{}
}
