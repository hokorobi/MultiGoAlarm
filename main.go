package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/rodolfoag/gow32"
)

func main() {
	_, err := gow32.CreateMutex("MultiGoAlarm")
	if err != nil {
		// TODO: 引数があったらアラームとして追加
		// fmt.Printf("Error: %d - %s\n", int(err.(syscall.Errno)), err.Error())
		os.Exit(0)
	}

	logfile, err := os.OpenFile(getLogFilename(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open logfile:" + err.Error())
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	log.SetFlags(log.Ldate | log.Ltime)

	app := newApp()

	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				// log.Println("tick")
				app.update()
			}
		}
		// t.Stop()
	}()

	ni := NotifyIcon(app.mw)
	defer ni.Dispose()

	if _, err := (declarative.MainWindow{
		AssignTo: &app.mw,
		Title:    "MultiGoAlarm",
		MinSize:  declarative.Size{Width: 400, Height: 300},
		Size:     declarative.Size{Width: 400, Height: 300},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.LineEdit{
						AssignTo: &app.time,
					},
					declarative.PushButton{
						Text:      "&Add",
						OnClicked: app.clickAdd,
					},
					declarative.PushButton{
						Text:      "&Quit",
						OnClicked: func() { app.mw.Close() },
					},
				},
			},
			declarative.ListBox{
				AssignTo:        &app.lb,
				Model:           app.list,
				OnItemActivated: app.lbItemActivated,
				Row:             10,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}

}

func getLogFilename() string {
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), filepath.Base(exec)+".log")
}

// App はこのアプリ全体の型
type app struct {
	mw   *walk.MainWindow
	time *walk.LineEdit
	lb   *walk.ListBox
	list *AlarmList
}

func newApp() app {
	var app app
	var err error
	app.mw, err = walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}
	app.list = NewAlarmList()
	return app
}

func (app *app) lbItemActivated() {
	if app.lb.CurrentIndex() < 0 {
		return
	}

	app.list.del(app.lb.CurrentIndex())
	app.lb.SetModel(app.list)
}
func (app *app) clickAdd() {
	item := NewAlarmItem(app.time.Text())
	if item == nil {
		walk.MsgBox(app.mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
		return
	}
	// debug
	// walk.MsgBox(mw, "confirm", item.start.String()+item.end.String()+item.message, walk.MsgBoxOK)
	app.list.add(*item)
	app.lb.SetModel(app.list)
}
func (app *app) update() {
	if len(app.list.list) < 1 {
		return
	}

	// log.Println("update")
	items := app.list.update()
	app.lb.SetModel(app.list)
	app.alarm(items)
}
func (app *app) alarm(items []AlarmItem) {
	for i := range items {
		go AlarmWindow(items[i].Message)
		time.Sleep(100 * time.Millisecond)
	}
}
