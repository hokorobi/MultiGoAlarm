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

type AlarmListJ struct {
	List []AlarmItem `json:"list"`
}

func NewAlarmListFile() AlarmListFile {
	var file AlarmListFile
	file.name = file.getFilename()
	file.mtime = file.getMtime()

	return file
}

func (file *AlarmListFile) getFilename() string {
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), filepath.Base(exec)+".json")
}

func (file *AlarmListFile) write(list *AlarmList) {
	var d AlarmListJ
	d.List = list.list
	b, err := json.MarshalIndent(&d, "", "  ")
	if err != nil {
		log.Println(err)
	}

	ioutil.WriteFile(file.name, b, os.ModePerm)
}

func (file *AlarmListFile) load(list *AlarmList) {
	if file.mtime == file.getMtime() {
		return
	}

	f, err := os.Open(file.name)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var d AlarmListJ
	err = dec.Decode(&d)
	if err != nil {
		log.Println(err)
		return
	}

	list.list = d.List
}

func (file *AlarmListFile) getMtime() time.Time {
	f1, err1 := os.Stat(file.name)
	if err1 == nil {
		return f1.ModTime()
	}
	// ファイルが存在しなかったら作成して変更時間を返す
	if !os.IsNotExist(err1) {
		log.Fatal(err1)
	}
	file.touch()
	f2, err2 := os.Stat(file.name)
	if err2 != nil {
		log.Fatal(err2)
	}
	return f2.ModTime()
}

func (file *AlarmListFile) touch() {
	f, err := os.Create(file.name)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}
