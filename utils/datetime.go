package utils

import "time"

func GetNowAsUnixTimestamp() int64 {
	now := time.Now()
	return now.Unix()
}
