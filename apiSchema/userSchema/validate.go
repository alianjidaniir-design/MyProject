package userSchema

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/statics/constants/status"
	"MyProject/statics/customErr"
	"fmt"
	"strings"
)

func (req *LoginRequest) Validate(validateExtraData commonSchema.ValidateExtraData) (string, int, error) {
	fmt.Println(req)
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	req.Family = strings.TrimSpace(req.Family)
	if req.Code == "" {
		return "03", status.StatusBadRequest, customErr.InvalidCode
	}
	if req.Name == "" {
		return "06", status.StatusBadRequest, customErr.InvalidName
	}
	if req.Family == "" {
		return "09", status.StatusBadRequest, customErr.InvalidFamily
	}
	_ = validateExtraData
	return "", status.StatusOK, nil
}

func (req *ListRequest) Validate(validateExtraData commonSchema.ValidateExtraData) (string, int, error) {
	_ = validateExtraData
	return "", status.StatusOK, nil
}
