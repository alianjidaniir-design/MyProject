package commonSchema

type BaseRequest[T any] struct {
	Body T
	Req  T
}

type ValidateExtraData struct {
	Headers map[string]string
}
