package util

import (
	"strconv"
	"time"
)

func GetCurrentTime() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format(time.RFC3339)
}

func GetCurrentTimeEx(timestamp string) string {
	tm := time.Now()
	inputTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		panic(err)
	}

	if !tm.After(inputTime) {
		tm = inputTime.Add(1 * time.Millisecond)
	}

	return tm.Format("2006-01-02T15:04:05.999Z07:00")
}

func GetCurrentUnixTime() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func IsTokenExpired(createdTime string, expiresIn int) bool {
	createdTimeObj, _ := time.Parse(time.RFC3339, createdTime)
	expiresAtObj := createdTimeObj.Add(time.Duration(expiresIn) * time.Second)
	return time.Now().After(expiresAtObj)
}
