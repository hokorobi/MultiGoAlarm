package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/lxn/walk"
)

type AlarmItems struct {
	walk.ListModelBase
	items []AlarmItem
}

func NewAlarmModel() *AlarmItems {
	m := &AlarmItems{items: make([]AlarmItem, 0)}
	return m
}

func (items *AlarmItems) add(item AlarmItem) {
	items.items = append(items.items, item)
	// items.write()
	return
}

func (items *AlarmItems) del(i int) {
	items.items = append(items.items[:i], items.items[i+1:]...)
}

func (items *AlarmItems) delId(id string) {
	for i := range items.items {
		if items.items[i].id == id {
			items.del(i)
			return
		}
	}
}

func (items *AlarmItems) update() []AlarmItem {
	var candidateItems []AlarmItem
	var candidateIds []string

	now := time.Now()
	for i := 0; i < len(items.items); i++ {
		// 終了時刻を過ぎている or 同じ
		if !items.items[i].end.After(now) {
			candidateItems = append(candidateItems, items.items[i])
			candidateIds = append(candidateIds, items.items[i].id)
		} else {
			items.items[i].setValue(now)
		}
	}
	for i := range candidateIds {
		items.delId(candidateIds[i])
	}

	return candidateItems
}

func (items *AlarmItems) write() {
	f, err := os.Create("timerlist.json")
	defer f.Close()
	if err != nil {
		log.Println(err)
		return
	}
	enc := json.NewEncoder(f)
	if err != nil {
		log.Println(err)
		return
	}
	err = enc.Encode(items)
	if err != nil {
		log.Println(err)
		return
	}
}

func (m *AlarmItems) ItemCount() int {
	return len(m.items)
}

func (m *AlarmItems) Value(index int) interface{} {
	return m.items[index].value
}
