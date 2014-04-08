package main

import (
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
	time    *walk.LineEdit
	message *walk.LineEdit
	list    *walk.ListBox
}

type SubWindow struct {
	*walk.MainWindow
}

type AlarmItem struct {
	start   *time.Time
	end     *time.Time
	message string
}

type AlarmItems struct {
	items []AlarmItem
}

func (items *AlarmItems) add(item AlarmItem) {
	items.items = append(items.items, item)
	return
}

func NewAlarmItem(timeString string, message string) *AlarmItem {
	start, end := GetTime(timeString)
	if start == nil {
		return nil
	}
	item := AlarmItem{start, end, message}
	return &item
}

func main() {

	ni := notifyIcon()
	defer ni.Dispose()

	mw := &MyMainWindow{}

	items := AlarmItems{}

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
						OnClicked: mw.clickAdd(items),
					},
					PushButton{
						Text:      "&Quit",
						OnClicked: mw.clickQuit,
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &mw.message,
					},
					ListBox{
						AssignTo: &mw.list,
						Row:      10,
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

func (mw *MyMainWindow) clickAdd(items AlarmItems) {
	item := NewAlarmItem(mw.time.Text(), mw.message.Text())
	if item == nil {
		walk.MsgBox(mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
		return
	}
	// debug
	walk.MsgBox(mw, "confirm", item.start.String()+item.end.String()+item.message, walk.MsgBoxOK)
	items.add(item)
}

func GetTime(s string) (*time.Time, *time.Time) {
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
	if result := re.MatchString(s); result {
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
