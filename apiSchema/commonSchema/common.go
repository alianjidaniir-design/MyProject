package commonSchema

type BaseRequest[T any] struct {
	Body    T                 `json:"body" msgpack:"body"`
	Headers map[string]string `json:"headers,omitempty" msgpack:"headers"`
}

type ValidateExtraData struct {
	Headers map[string]string
}
