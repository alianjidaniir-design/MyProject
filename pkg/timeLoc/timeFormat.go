package timeLoc

import (
	"errors"
	"time"
)

func FormatTime(input string) (string, error) {
	// سعی کن رشته را به ساعت تبدیل کنی
	t, err := time.Parse("15:04", input)
	if err != nil {
		return "", errors.New("فرمت ساعت نامعتبر است")
	}
	// خروجی استاندارد دیتابیس
	return t.Format("15:04:05"), nil
}
