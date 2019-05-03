package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/lxn/walk"
)

type AlarmItems struct {
	filename string
	walk.ListModelBase
	items []AlarmItem
}

type AlarmItemsJ struct {
	Items []AlarmItem `json:"items"`
}

func NewAlarmModel() *AlarmItems {
	m := &AlarmItems{items: make([]AlarmItem, 0)}
	m.filename = m.getAlarmsFilename()
	return m
}

func (items *AlarmItems) add(item AlarmItem) {
	items.items = append(items.items, item)
	items.write()
}

func (items *AlarmItems) del(i int) {
	items.items = append(items.items[:i], items.items[i+1:]...)
	items.write()
}

func (items *AlarmItems) delID(id string) {
	for i := range items.items {
		if items.items[i].ID == id {
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
		if items.items[i].isTimeUp(now) {
			candidateItems = append(candidateItems, items.items[i])
			candidateIds = append(candidateIds, items.items[i].ID)
		} else {
			items.items[i].setValue(now)
		}
	}
	for i := range candidateIds {
		items.delID(candidateIds[i])
	}

	return candidateItems
}

func (items *AlarmItems) load() {
	var d AlarmItemsJ

	f, err := os.Open(items.filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&d)
	if err != nil {
		log.Println(err)
		return
	}

	items.items = d.Items
}

func (items *AlarmItems) write() {
	var d AlarmItemsJ
	d.Items = items.items
	b, err := json.MarshalIndent(&d, "", "  ")
	if err != nil {
		log.Println(err)
	}

	ioutil.WriteFile(items.filename, b, os.ModePerm)
}

func (items *AlarmItems) getAlarmsFilename() string {
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), filepath.Base(exec)+".json")
}

func (items *AlarmItems) ItemCount() int {
	return len(items.items)
}

func (items *AlarmItems) Value(index int) interface{} {
	return items.items[index].Value
}
