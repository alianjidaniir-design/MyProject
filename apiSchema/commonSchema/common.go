package commonSchema

type BaseRequest[T any] struct {
	Body    T                 `json:"body"`
	Headers map[string]string `json:"headers,omitempty"`
}

type ValidateExtraData struct {
	Headers map[string]string
}
