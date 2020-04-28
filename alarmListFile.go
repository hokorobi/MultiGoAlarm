package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

// AlarmListFile はアラームのリスト用のファイルを操作する型
type AlarmListFile struct {
	name  string
	mtime time.Time
}

type alarmListForJSON struct {
	List []AlarmItem `json:"list"`
}

// NewAlarmListFile は AlarmListFile を新規作成する関数
func NewAlarmListFile() AlarmListFile {
	var file AlarmListFile
	file.name = getFilename(".json")

	return file
}

func (file *AlarmListFile) write(list *AlarmList) {
	var d alarmListForJSON
	d.List = list.list
	b, err := json.MarshalIndent(&d, "", "  ")
	if err != nil {
		Logg(err)
	}

	ioutil.WriteFile(file.name, b, os.ModePerm)
	file.mtime = file.getMtime()
}

func (file *AlarmListFile) load(list *AlarmList) {
	if !file.mtime.IsZero() && file.mtime == file.getMtime() {
		return
	}

	f, err := os.Open(file.name)
	if err != nil {
		Logg(err)
		return
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var d alarmListForJSON
	err = dec.Decode(&d)
	if err != nil {
		Logg(err)
		return
	}

	file.mtime = file.getMtime()
	list.list = d.List
}

func (file *AlarmListFile) getMtime() time.Time {
	f1, err1 := os.Stat(file.name)
	if err1 == nil {
		return f1.ModTime()
	}
	// ファイルが存在しなかったら作成して変更時間を返す
	if !os.IsNotExist(err1) {
		Logf(err1)
	}
	file.write(&AlarmList{list: make([]AlarmItem, 0)})
	f2, err2 := os.Stat(file.name)
	if err2 != nil {
		Logf(err2)
	}
	return f2.ModTime()
}
