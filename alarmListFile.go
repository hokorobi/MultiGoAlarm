package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/hokorobi/go-utils/logutil"
)

type alarmListFile struct {
	name  string
	mtime time.Time
}

type alarmListForJSON struct {
	List []alarmItem `json:"list"`
}

func newAlarmListFile() alarmListFile {
	var file alarmListFile
	file.name = getFilename(".json")

	return file
}

func (file *alarmListFile) write(list *alarmList) {
	var d alarmListForJSON
	d.List = list.list
	b, err := json.MarshalIndent(&d, "", "  ")
	if err != nil {
		logutil.PrintTee(err)
	}

	ioutil.WriteFile(file.name, b, os.ModePerm)
	file.mtime = file.getMtime()
}

func (file *alarmListFile) load(list *alarmList) {
	if !file.mtime.IsZero() && file.mtime == file.getMtime() {
		return
	}

	f, err := os.Open(file.name)
	if os.IsNotExist(err) {
		// ファイルがなければ何もせずに新規作成
		return
	} else if err != nil {
		logutil.PrintTee(err)
		return
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var d alarmListForJSON
	err = dec.Decode(&d)
	if err != nil {
		logutil.PrintTee(err)
		return
	}

	file.mtime = file.getMtime()
	list.list = d.List
}

func (file *alarmListFile) getMtime() time.Time {
	f1, err1 := os.Stat(file.name)
	if err1 == nil {
		return f1.ModTime()
	}
	// ファイルが存在しなかったら作成して変更時間を返す
	if !os.IsNotExist(err1) {
		logutil.FatalTee(err1)
	}
	file.write(&alarmList{list: make([]alarmItem, 0)})
	f2, err2 := os.Stat(file.name)
	if err2 != nil {
		logutil.FatalTee(err2)
	}
	return f2.ModTime()
}
