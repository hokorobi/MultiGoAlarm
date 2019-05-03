package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type AlarmListFile struct {
	name  string
	mtime time.Time
}

type AlarmItemsJ struct {
	List []AlarmItem `json:"list"`
}

func NewAlarmListFile() AlarmListFile {
	var f AlarmListFile
	f.name = f.getFilename()
	return f
}

func (file *AlarmListFile) getFilename() string {
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), filepath.Base(exec)+".json")
}

func (file *AlarmListFile) write(list *AlarmList) {
	var d AlarmItemsJ
	d.List = list.list
	b, err := json.MarshalIndent(&d, "", "  ")
	if err != nil {
		log.Println(err)
	}

	ioutil.WriteFile(file.name, b, os.ModePerm)
}

func (file *AlarmListFile) load(list *AlarmList) {
	var d AlarmItemsJ

	f, err := os.Open(file.name)
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
