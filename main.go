package main

import (
	"log"
	"os"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/rodolfoag/gow32"
)

func main() {
	_, err := gow32.CreateMutex("MultiGoAlarm")
	if err != nil {
		// fmt.Printf("Error: %d - %s\n", int(err.(syscall.Errno)), err.Error())
		os.Exit(0)
	}
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
					go AlarmWindow(alarmItems[i].Message)
					time.Sleep(100 * time.Millisecond)
				}
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
