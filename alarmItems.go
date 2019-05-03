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

type AlarmList struct {
	filename string
	walk.ListModelBase
	list []AlarmItem
}

type AlarmItemsJ struct {
	List []AlarmItem `json:"list"`
}

func NewAlarmList() *AlarmList {
	m := &AlarmList{list: make([]AlarmItem, 0)}
	m.filename = m.getAlarmsFilename()
	return m
}

func (list *AlarmList) add(item AlarmItem) {
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
	var d AlarmItemsJ

	f, err := os.Open(list.filename)
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

	list.list = d.List
}

func (list *AlarmList) write() {
	var d AlarmItemsJ
	d.List = list.list
	b, err := json.MarshalIndent(&d, "", "  ")
	if err != nil {
		log.Println(err)
	}

	ioutil.WriteFile(list.filename, b, os.ModePerm)
}

func (list *AlarmList) getAlarmsFilename() string {
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), filepath.Base(exec)+".json")
}

func (list *AlarmList) ItemCount() int {
	return len(list.list)
}

func (list *AlarmList) Value(index int) interface{} {
	return list.list[index].Value
}
