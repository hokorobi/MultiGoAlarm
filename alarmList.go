package main

import (
	"time"

	"github.com/lxn/walk"
)

type AlarmList struct {
	file AlarmListFile
	walk.ListModelBase
	list []AlarmItem
}

// NewAlarmList は AlarmList を生成する関数
func NewAlarmList() *AlarmList {
	m := &AlarmList{list: make([]AlarmItem, 0)}
	m.file = NewAlarmListFile()
	m.load()
	m.deleteTimeout()
	return m
}

func (list *AlarmList) add(item AlarmItem) {
	list.load()
	list.list = append(list.list, item)
	list.write()
}

func (list *AlarmList) del(i int) {
	list.list = append(list.list[:i], list.list[i+1:]...)
	list.write()
}

func (list *AlarmList) delID(id string) {
	for i := range list.list {
		if list.list[i].ID == id {
			list.del(i)
			return
		}
	}
}

func (list *AlarmList) update() []AlarmItem {
	var candidateItems []AlarmItem
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

func (list *AlarmList) load() {
	list.file.load(list)
}

func (list *AlarmList) write() {
	list.file.write(list)
}

func (list *AlarmList) ItemCount() int {
	return len(list.list)
}

func (list *AlarmList) Value(index int) interface{} {
	return list.list[index].Value
}

func (list *AlarmList) deleteTimeout() {
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
