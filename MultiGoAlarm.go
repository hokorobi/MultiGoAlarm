package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MyMainWindow struct {
	*walk.MainWindow
	time  *walk.LineEdit
	lb    *walk.ListBox
	model *AlarmItems
}

type SubWindow struct {
	*walk.MainWindow
}

type AlarmItem struct {
	start   *time.Time
	end     *time.Time
	message string
	value   string
}

type AlarmItems struct {
	walk.ListModelBase
	items []AlarmItem
}

func (item *AlarmItem) getTime(s string) (*time.Time, *time.Time) {
	start := time.Now()
	// 数字だけなら分として扱う
	if d, err := time.ParseDuration(s + "m"); err == nil {
		end := start.Add(d)
		return &start, &end
	}
	// 1h2m などを解釈
	if d, err := time.ParseDuration(s); err == nil {
		end := start.Add(d)
		return &start, &end
	}
	// hh:mm
	re := regexp.MustCompile("^[0-9]+:[0-9]+$")
	if re.MatchString(s) {
		hhmm := strings.Split(s, ":")
		hh, _ := strconv.Atoi(hhmm[0])
		mm, _ := strconv.Atoi(hhmm[1])
		end := time.Date(start.Year(), start.Month(), start.Day(), hh, mm, 0, 0, start.Location())
		// 翌日の hh:mm
		if start.After(end) {
			end = end.Add(time.Hour * 24)
		}
		return &start, &end
	}

	return nil, nil
}

func (items *AlarmItems) add(item AlarmItem) {
	items.items = append(items.items, item)
	// items.write()
	return
}

func (items *AlarmItems) del(i int) {
	items.items = append(items.items[:i], items.items[i+1:]...)
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

func NewAlarmItem(s string) *AlarmItem {
	var message string
	var timeString string

	if strings.Index(s, " ") > 0 {
		timeString = s[0:strings.Index(s, " ")]
		message = s[strings.Index(s, " "):]
	} else {
		timeString = s
		message = ""
	}

	item := new(AlarmItem)
	start, end := item.getTime(timeString)
	if start == nil {
		return nil
	}
	item.start = start
	item.end = end
	item.message = message
	item.value = end.Sub(*start).String() + " " + message

	return item
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
	// t := time.NewTicker(3 * time.Second)
	// for {
	// 	select {
	// 	case <-t.C:
	// 		mw.updatelist()
	// 	}
	// 	t.Stop()
	// }
}

func (mw *MyMainWindow) lb_ItemActivated() {
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

func (mw *MyMainWindow) updatelist() {
	return
}

func NewAlarmModel() *AlarmItems {
	m := &AlarmItems{items: make([]AlarmItem, 0)}
	return m
}

func (m *AlarmItems) ItemCount() int {
	return len(m.items)
}

func (m *AlarmItems) Value(index int) interface{} {
	return m.items[index].value
}

func Alarm() {
	sw := &SubWindow{}

	if _, err := (MainWindow{
		AssignTo: &sw.MainWindow,
		Title:    "Alarm",
		MinSize:  Size{300, 300},
		Layout:   VBox{},
		Children: []Widget{
			Label{
				Text: "Name:",
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
