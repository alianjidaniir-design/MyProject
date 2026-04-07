package customErr

import (
	"MyProject/statics/constants"
	"errors"
)

var (
	BadRequest        = errors.New(constants.BadRequest)
	ServerError       = errors.New(constants.ServerError)
	Forbidden         = errors.New(constants.Forbidden)
	Success           = errors.New(constants.Success)
	Error             = errors.New(constants.Error)
	RequestOK         = errors.New(constants.RequestOK)
	OK                = errors.New(constants.OK)
	UnAuthorized      = errors.New(constants.UnAuthorized)
	InvalidCode       = errors.New(constants.InvalidCode)
	InvalidName       = errors.New(constants.InvalidName)
	InvalidFamily     = errors.New(constants.InvalidFamily)
	InvalidTitle      = errors.New(constants.InvalidTitle)
	InvalidCourseCode = errors.New(constants.InvalidCourseCode)
)
