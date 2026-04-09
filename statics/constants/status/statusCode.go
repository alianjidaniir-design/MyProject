package status

const (
	StatusOK           = 200
	StatusBadRequest   = 400
	StatusUnauthorized = 401
	StatusForbidden    = 403
	StatusNotFound     = 404
	
	TooManyRequests           = 429
	StatusInternalServerError = 500
	UnAvailableServiceError   = 503
	RedirectPermanently       = 301
)
