package commonSchema

type BaseRequest[T any] struct {
	Body    T
	Headers map[string]string
}

type ValidateExtraData struct {
	Headers map[string]string
}
