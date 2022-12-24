package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/hokorobi/go-utils/logutil"
	"github.com/lxn/win"
	"github.com/rodolfoag/gow32"
)

type app struct {
	list *alarmList
}

func main() {
	var (
		listAlarms = flag.Bool("l", false, "Print timer list.")
	)
	flag.Parse()

	if *listAlarms {
		var templist = loadAlarmList()
		templist.sort()
		var message = ""
		for i := range templist.list {
			var item = templist.list[i]
			if message != "" {
				message = message + "\n"
			}
			message = message + item.End.Format("15:04:05") + " " + item.Message
		}
		messageBox(message, "", win.MB_OK)
	}

	if len(flag.Args()) > 0 {
		item := newAlarmItem(strings.Join(flag.Args(), " "))
		if item == nil {
			messageBox(
				"Error: Enter valid time format:"+strings.Join(flag.Args(), " "),
				"Error",
				win.MB_OK+win.MB_ICONEXCLAMATION)
		} else {
			templist := loadAlarmList()
			templist.add(*item)
		}
	}

	_, err := gow32.CreateMutex("MultiGoAlarm")
	if err != nil {
		os.Exit(0)
	}

	logutil.PrintTee("Start")
	defer logutil.PrintTee("End")

	app := newApp()

	sc := gocron.NewScheduler(time.UTC)
	sc.Every(1).Seconds().Do(app.update)
	sc.StartBlocking()
}

func newApp() app {
	var app app
	app.list = loadAlarmList()
	return app
}
func (app *app) update() {
	items := app.list.update()
	app.alarm(items)
}
func (app *app) alarm(items []alarmItem) {
	for i := range items {
		go alarm(items[i].Message)
		logutil.PrintTee("Alarm: " + items[i].End.Format("15:04:05") + " " + items[i].Message)
		time.Sleep(100 * time.Millisecond)
	}
}

func messageBox(message, title string, uType uint32) {
	win.MessageBox(
		win.HWND(0),
		UTF16PtrFromString(message),
		UTF16PtrFromString(title),
		uType)
}

// "Go から Windows の MessageBox を呼び出す - Qiita" https://qiita.com/manymanyuni/items/867d7e0112ce22dec6d5
func UTF16PtrFromString(s string) *uint16 {
	result, _ := syscall.UTF16PtrFromString(s)
	return result
}

// https://qiita.com/KemoKemo/items/d135ddc93e6f87008521#comment-7d090bd8afe54df429b9
func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
func getFilename(ext string) string {
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), getFileNameWithoutExt(exec)+ext)
}
