package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

var lock = sync.RWMutex{}

func main() {
	logfile, err := os.OpenFile("./test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open test.log:" + err.Error())
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	log.SetFlags(log.Ldate | log.Ltime)

	var app App
	app.mw, err = walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	app.model = NewAlarmModel()

	var alarmItems []AlarmItem

	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				// log.Println("tick")
				alarmItems = app.update()
				for i := range alarmItems {
					go AlarmWindow(alarmItems[i].message)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
		// t.Stop()
	}()

	ni := notifyIcon(app.mw)
	defer ni.Dispose()

	if _, err := (declarative.MainWindow{
		AssignTo: &app.mw,
		Title:    "MultiGoAlarm",
		MinSize:  declarative.Size{Width: 400, Height: 300},
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
						OnClicked: app.clickQuit,
					},
				},
			},
			declarative.ListBox{
				AssignTo:        &app.lb,
				Model:           app.model,
				OnItemActivated: app.lbItemActivated,
				Row:             10,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}

}

type App struct {
	mw    *walk.MainWindow
	time  *walk.LineEdit
	lb    *walk.ListBox
	model *AlarmItems
}

func (app *App) lbItemActivated() {
	if app.lb.CurrentIndex() < 0 {
		return
	}

	app.model.del(app.lb.CurrentIndex())
	app.lb.SetModel(app.model)
}

func (app *App) clickAdd() {
	item := NewAlarmItem(app.time.Text())
	if item == nil {
		walk.MsgBox(app.mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
		return
	}
	// debug
	// walk.MsgBox(mw, "confirm", item.start.String()+item.end.String()+item.message, walk.MsgBoxOK)
	app.model.add(*item)
	app.lb.SetModel(app.model)
}

func (app *App) update() []AlarmItem {
	if len(app.model.items) <= 0 {
		return nil
	}

	// log.Println("update")
	items := app.model.update()
	app.lb.SetModel(app.model)
	return items
}

func (app *App) clickQuit() {
	os.Exit(0)
}

func notifyIcon(mw *walk.MainWindow) *walk.NotifyIcon {
	// load icon
	icon, err := walk.NewIconFromFile("MultiGoAlarm.ico")
	if err != nil {
		log.Fatal(err)
	}

	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}

	// Set the icon and a tool tip text.
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal(err)
	}
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}
	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}
	return ni
}
