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

	mw := &MyMainWindow{model: NewAlarmModel()}

	var alarmItems []AlarmItem

	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				// log.Println("tick")
				alarmItems = mw.update()
				for i := range alarmItems {
					go Alarm(alarmItems[i].message)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
		// t.Stop()
	}()

	// FIXME: notifyIcon() 内でメインウィンドウを作っているので、こいつが終了したら通常のメインウィンドウが出てくる
	// notifyIcon()

	if _, err := (declarative.MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "MultiGoAlarm",
		MinSize:  declarative.Size{Width: 400, Height: 300},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.LineEdit{
						AssignTo: &mw.time,
					},
					declarative.PushButton{
						Text:      "&Add",
						OnClicked: mw.clickAdd,
					},
					declarative.PushButton{
						Text:      "&Quit",
						OnClicked: mw.clickQuit,
					},
				},
			},
			declarative.ListBox{
				AssignTo:        &mw.lb,
				Model:           mw.model,
				OnItemActivated: mw.lbItemActivated,
				Row:             10,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}

}

type MyMainWindow struct {
	*walk.MainWindow
	time  *walk.LineEdit
	lb    *walk.ListBox
	model *AlarmItems
}

func (mw *MyMainWindow) lbItemActivated() {
	if mw.lb.CurrentIndex() < 0 {
		return
	}

	mw.model.del(mw.lb.CurrentIndex())
	mw.lb.SetModel(mw.model)
}

func (mw *MyMainWindow) clickAdd() {
	item := NewAlarmItem(mw.time.Text())
	if item == nil {
		walk.MsgBox(mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
		return
	}
	// debug
	// walk.MsgBox(mw, "confirm", item.start.String()+item.end.String()+item.message, walk.MsgBoxOK)
	mw.model.add(*item)
	mw.lb.SetModel(mw.model)
}

func (mw *MyMainWindow) update() []AlarmItem {
	if len(mw.model.items) <= 0 {
		return nil
	}

	// log.Println("update")
	items := mw.model.update()
	mw.lb.SetModel(mw.model)
	return items
}

func (mw *MyMainWindow) clickQuit() {
	os.Exit(0)
}

type SubWindow struct {
	*walk.MainWindow
}

func Alarm(s string) {
	var message string

	sw := &SubWindow{}

	lock.RLock()
	defer lock.RUnlock()

	if s == "" {
		message = "It is Time!"
	} else {
		message = s
	}

	if _, err := (declarative.MainWindow{
		AssignTo: &sw.MainWindow,
		Title:    "Alarm",
		MinSize:  declarative.Size{Width: 300, Height: 300},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Label{
				Text: message,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

func notifyIcon() {
	// load icon
	icon, err := walk.NewIconFromFile("MultiGoAlarm.ico")
	if err != nil {
		log.Fatal(err)
	}

	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()

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
	mw.Run()
}
