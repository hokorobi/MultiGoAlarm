package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-toast/toast"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/rodolfoag/gow32"
)

func main() {
	_, err := gow32.CreateMutex("MultiGoAlarm")
	if err != nil {
		if len(os.Args) > 1 {
			item := newAlarmItem(strings.Join(os.Args[1:], " "))
			if item == nil {
				logf("Error: Enter valid time format:" + strings.Join(os.Args[1:], " "))
			}
			list := newAlarmList()
			list.add(*item)
			notification(*item)
			os.Exit(0)
		}
		os.Exit(1)
	}

	app := newApp()

	t := time.NewTicker(time.Second)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				app.update()
			}
		}
	}()

	app.ni = newNotifyIcon(&app)
	defer app.ni.Dispose()

	if len(os.Args) > 1 {
		item := newAlarmItem(strings.Join(os.Args[1:], " "))
		if item == nil {
			walk.MsgBox(app.mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
		} else {
			app.list.add(*item)
		}
	}

	logg("Run.")
	defer logg("Stop.")

	if _, err := (declarative.MainWindow{
		AssignTo: &app.mw,
		Title:    "MultiGoAlarm",
		Visible:  false,
	}.Run()); err != nil {
		logf(err)
	}
}

// App はこのアプリ全体の型
type app struct {
	mw   *walk.MainWindow
	list *alarmList
	ni   *walk.NotifyIcon
}

func newApp() app {
	var app app
	var err error
	app.mw, err = walk.NewMainWindow()
	if err != nil {
		logf(err)
	}
	app.list = newAlarmList()
	return app
}
func (app *app) update() {
	items := app.list.update()
	app.alarm(items)
}
func (app *app) alarm(items []alarmItem) {
	for i := range items {
		go alarm(items[i].Message)
		logg("Alarm: " + items[i].End.Format("15:04:05") + " " + items[i].Message)
		time.Sleep(100 * time.Millisecond)
	}
}

func notification(item alarmItem) {
	iconpath, err := filepath.Abs("alarm-check.png")
	if err != nil {
		logg(err)
	}

	notify := toast.Notification{
		AppID:   "MultiGoAlarm",
		Title:   "Add Alarm",
		Icon:    iconpath,
		Message: item.End.Format("15:04:05") + " " + item.Message,
	}
	err = notify.Push()
	if err != nil {
		logg(err)
	}
}

// https://qiita.com/KemoKemo/items/d135ddc93e6f87008521#comment-7d090bd8afe54df429b9
func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
func getFilename(ext string) string {
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), getFileNameWithoutExt(exec)+ext)
}

func logg(m interface{}) {
	f, err := os.OpenFile(getFilename(".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("Cannot open log file: " + err.Error())
	}
	defer f.Close()

	log.SetOutput(io.MultiWriter(f, os.Stderr))
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println(m)
}
func logf(m interface{}) {
	logg(m)
	os.Exit(1)
}
