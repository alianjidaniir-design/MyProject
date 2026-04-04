package userSchema

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/statics/constants/status"
	"fmt"
	"strings"
)

func (req *LoginRequest) Validate(validateExtraData commonSchema.ValidateExtraData) (string, int, error) {
	fmt.Println(req)
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	req.Family = strings.TrimSpace(req.Family)

	_ = validateExtraData
	return "", status.StatusOK, nil
}

func (req *ListRequest) Validate(validateExtraData commonSchema.ValidateExtraData) (string, int, error) {
	_ = validateExtraData
	return "", status.StatusOK, nil
}
