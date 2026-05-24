package timeLoc

import "time"

func MyLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return location
}
