package oktools

import "time"

type (
	SortedTimeStringMapList []map[string]string
)

func (s SortedTimeStringMapList) Len() int {
	return len(s)
}

func (s SortedTimeStringMapList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedTimeStringMapList) Less(i, j int) bool {
	time1Str, ok1 := s[i]["title"]
	time2Str, ok2 := s[j]["title"]

	if !ok1 || !ok2 {
		return true
	}

	time1, err1 := time.ParseInLocation("15:04:05", time1Str, time.Local)
	time2, err2 := time.ParseInLocation("15:04:05", time2Str, time.Local)

	if err1 != nil || err2 != nil {
		return true
	}

	return time1.Before(time2) //s[i]["title"] < s[j]["title"]
}

func (s SortedTimeStringMapList) Get() []map[string]string {
	return s
}
