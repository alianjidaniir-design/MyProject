package timeLoc

import (
	"MyProject/statics/constants"
	"errors"
	"fmt"
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

func FormatDataTime(input string) (*time.Time, error) {

	in, err := time.ParseInLocation(time.DateTime, input, MyLocation())
	if err != nil {
		return &time.Time{}, err
	}
	timeIn := in.In(MyLocation())
	return &timeIn, nil
}

func CheckDuration(in string, out string) error {
	time1, err := time.Parse("15:04", in)
	if err != nil {
		return err
	}
	time2, err := time.Parse("15:04", out)
	if err != nil {
		return err
	}
	dur := time2.Sub(time1).Hours()
	if dur != constants.DurationTime {
		return fmt.Errorf("the duration of each class should be %.2f hour", constants.DurationTime)
	}
	return nil
}

func CheckTimeExam(in *time.Time, out *time.Time) error {
	if in.After(*out) {
		return fmt.Errorf("time finish should be next time start")
	}
	dur := out.Sub(*in).Hours()
	if dur > constants.DurationTime {
		return fmt.Errorf("exam Time is very Long")
	} else if dur == 0 {
		return fmt.Errorf("can not finish time like with start time")
	}
	return nil
}
