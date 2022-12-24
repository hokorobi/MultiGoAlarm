package main

import (
	"time"

	"github.com/go-toast/toast"
	"github.com/lxn/walk"
)

type alarmList struct {
	file alarmListFile
	walk.ListModelBase
	list []alarmItem
}

func loadAlarmList() *alarmList {
	m := &alarmList{list: make([]alarmItem, 0)}
	m.file = newAlarmListFile()
	m.load()
	m.deleteTimeout()
	return m
}

func (list *alarmList) add(item alarmItem) {
	list.load()
	list.list = append(list.list, item)
	list.write()
	logg("Add Alarm: " + item.End.Format("15:04:05") + " " + item.Message)
	notification(item)
}

func notification(item alarmItem) {
	notify := toast.Notification{
		AppID:   "MultiGoAlarm",
		Title:   "Add Alarm",
		Message: item.End.Format("15:04:05") + " " + item.Message,
	}
	err := notify.Push()
	if err != nil {
		logg(err)
	}
}

func (list *alarmList) del(i int) {
	list.list = append(list.list[:i], list.list[i+1:]...)
	list.write()
}

func (list *alarmList) delID(id string) {
	for i := range list.list {
		if list.list[i].ID == id {
			list.del(i)
			return
		}
	}
}

func (list *alarmList) update() []alarmItem {
	var candidateItems []alarmItem
	var candidateIds []string

	now := time.Now()
	list.load()
	for i := 0; i < len(list.list); i++ {
		// 終了時刻を過ぎている or 同じ
		if list.list[i].isTimeUp(now) {
			candidateItems = append(candidateItems, list.list[i])
			candidateIds = append(candidateIds, list.list[i].ID)
		} else {
			list.list[i].setValue(now)
		}
	}
	for i := range candidateIds {
		list.delID(candidateIds[i])
	}

	return candidateItems
}

func (list *alarmList) load() {
	list.file.load(list)
}

func (list *alarmList) write() {
	list.file.write(list)
}

func (list *alarmList) ItemCount() int {
	return len(list.list)
}

func (list *alarmList) Value(index int) interface{} {
	return list.list[index].Value
}

func (list *alarmList) deleteTimeout() {
	var isDelete bool
	for _, e := range list.list {
		if e.End.Before(time.Now()) {
			list.delID(e.ID)
			isDelete = true
		}
	}
	if isDelete {
		list.write()
	}
}
