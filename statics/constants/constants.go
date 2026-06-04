package constants

import "time"

const (
	MaxNumberUnits = 9
	PageSize       = 50
)
const (
	AccessTokenExpiry  = 20 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

const (
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
)

const (
	StatusUnEnrolled   = "unenrolled"
	StatusEnrolled     = "enrolled"
	StatusCanceled     = "canceled"
	StatusReserveation = "reserveation"
)

const (
	Passed = "Passed"
	Failed = "failed"
)

const (
	Approved = "approved"
	Rejected = "rejected"
	Review   = "review"
)

const (
	Installment = "installment"
	Cash        = "cash"
)

const (
	Melat   = "melat"
	Meli    = "meli"
	Saderat = "saderat"
)

const (
	JSON    = "json"
	MSGPACK = "msgpack"
)

const (
	Success           = "success"
	Error             = "error"
	RequestOK         = "requestOk"
	OK                = "ok"
	BadRequest        = "badRequest"
	ServerError       = "serverError"
	Forbidden         = "forbidden"
	UnAuthorized      = "unAuthorized"
	InvalidCode       = "invalidCode"
	InvalidName       = "invalidName"
	InvalidFamily     = "invalidFamily"
	InvalidTitle      = "invalidTitle"
	InvalidCourseCode = "invalidCourseCode"
)
