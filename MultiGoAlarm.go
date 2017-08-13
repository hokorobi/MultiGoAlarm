package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type AlarmItems struct {
	walk.ListModelBase
	items []AlarmItem
}

func NewAlarmModel() *AlarmItems {
	m := &AlarmItems{items: make([]AlarmItem, 0)}
	return m
}

func (items *AlarmItems) add(item AlarmItem) {
	items.items = append(items.items, item)
	// items.write()
	return
}

func (items *AlarmItems) del(i int) {
	items.items = append(items.items[:i], items.items[i+1:]...)
}

func (items *AlarmItems) delId(id string) {
	for i := range items.items {
		if items.items[i].id == id {
			items.del(i)
			return
		}
	}
}

func (items *AlarmItems) update() []AlarmItem {
	var candidateItems []AlarmItem
	var candidateIds []string

	now := time.Now()
	for i := 0; i < len(items.items); i++ {
		// 終了時刻を過ぎている or 同じ
		if !items.items[i].end.After(now) {
			candidateItems = append(candidateItems, items.items[i])
			candidateIds = append(candidateIds, items.items[i].id)
		} else {
			items.items[i].setValue(now)
		}
	}
	for i := range candidateIds {
		items.delId(candidateIds[i])
	}

	return candidateItems
}

func (items *AlarmItems) write() {
	f, err := os.Create("timerlist.json")
	defer f.Close()
	if err != nil {
		log.Println(err)
		return
	}
	enc := json.NewEncoder(f)
	if err != nil {
		log.Println(err)
		return
	}
	err = enc.Encode(items)
	if err != nil {
		log.Println(err)
		return
	}
}

func (m *AlarmItems) ItemCount() int {
	return len(m.items)
}

func (m *AlarmItems) Value(index int) interface{} {
	return m.items[index].value
}

type MyMainWindow struct {
	*walk.MainWindow
	time  *walk.LineEdit
	lb    *walk.ListBox
	model *AlarmItems
}

type SubWindow struct {
	*walk.MainWindow
}

func main() {
	logfile, err := os.OpenFile("./test.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open test.log:" + err.Error())
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	log.SetFlags(log.Ldate | log.Ltime)

	ni := notifyIcon()
	defer ni.Dispose()

	mw := &MyMainWindow{model: NewAlarmModel()}

	var items []AlarmItem

	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				// log.Println("tick")
				items = mw.update()
				// for i := range items {
				// 	Alarm(items[i].message)
				// }
			}
		}
		t.Stop()
	}()

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "MultiGoAlarm",
		MinSize:  Size{400, 300},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &mw.time,
					},
					PushButton{
						Text:      "&Add",
						OnClicked: mw.clickAdd,
					},
					PushButton{
						Text:      "&Quit",
						OnClicked: mw.clickQuit,
					},
				},
			},
			ListBox{
				AssignTo:        &mw.lb,
				Model:           mw.model,
				OnItemActivated: mw.lb_ItemActivated,
				Row:             10,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}

}

func (mw *MyMainWindow) lb_ItemActivated() {
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

func Alarm(message string) {
	sw := &SubWindow{}

	if _, err := (MainWindow{
		AssignTo: &sw.MainWindow,
		Title:    "Alarm",
		MinSize:  Size{300, 300},
		Layout:   VBox{},
		Children: []Widget{
			Label{
				Text: message,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

func (mw *MyMainWindow) clickQuit() {
	os.Exit(0)
}

func notifyIcon() *walk.NotifyIcon {
	// load icon
	icon, err := walk.NewIconFromFile("MultiGoAlarm.ico")
	if err != nil {
		log.Fatal(err)
	}

	ni, err := walk.NewNotifyIcon()
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
