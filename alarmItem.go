package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AlarmItem struct {
	start   *time.Time
	end     *time.Time
	message string
	value   string
	id      string
}

func (item *AlarmItem) setValue(start time.Time) {
	hour := "00"
	minute := "00"
	second := "00"
	var index int

	v := item.end.Sub(start).String()
	index = strings.Index(v, "h")
	if index > -1 {
		hour = v[:index]
		v = v[index+1:]
	}
	index = strings.Index(v, "m")
	if index > -1 {
		minute = v[:index]
		v = v[index+1:]
	}
	index = strings.Index(v, ".")
	if index > -1 {
		second = v[:index]
	} else {
		second = v[:strings.Index(v, "s")]
	}
	item.value = fmt.Sprintf("%02s:%02s:%02s %s", hour, minute, second, item.message)
}

func (item *AlarmItem) getTime(s string) (*time.Time, *time.Time) {
	start := time.Now()
	// 数字だけなら分として扱う
	if d, err := time.ParseDuration(s + "m"); err == nil {
		end := start.Add(d)
		return &start, &end
	}
	// 1h2m などを解釈
	if d, err := time.ParseDuration(s); err == nil {
		end := start.Add(d)
		return &start, &end
	}
	// hh:mm
	re := regexp.MustCompile("^[0-9]+:[0-9]+$")
	if re.MatchString(s) {
		hhmm := strings.Split(s, ":")
		hh, _ := strconv.Atoi(hhmm[0])
		mm, _ := strconv.Atoi(hhmm[1])
		end := time.Date(start.Year(), start.Month(), start.Day(), hh, mm, 0, 0, start.Location())
		// 翌日の hh:mm
		if start.After(end) {
			end = end.Add(time.Hour * 24)
		}
		return &start, &end
	}

	return nil, nil
}

func NewAlarmItem(s string) *AlarmItem {
	var message string
	var timeString string

	if strings.Index(s, " ") > 0 {
		timeString = s[0:strings.Index(s, " ")]
		message = s[strings.Index(s, " "):]
	} else {
		timeString = s
		message = ""
	}

	item := new(AlarmItem)
	start, end := item.getTime(timeString)
	if start == nil {
		return nil
	}
	item.start = start
	item.end = end
	item.message = message
	item.setValue(*start)
	item.id = fmt.Sprint(start.Unix())

	return item
}

func (item *AlarmItem) isTimeUp(now time.Time) bool {
	if item.end.Sub(now).Seconds() < 1.0 {
		return true
	} else {
		return false
	}
}
