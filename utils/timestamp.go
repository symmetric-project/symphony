package utils

import "time"

func CurrentTimestamp() int {
	return int(time.Now().Unix())
}
