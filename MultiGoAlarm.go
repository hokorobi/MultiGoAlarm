package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	// "strings"
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

func (mw *MyMainWindow) clickAdd() {
	item := NewAlarmItem(mw.time.Text(), mw.message.Text())
	if item == nil {
		walk.MsgBox(mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
		return
	}
	// debug
	walk.MsgBox(mw, "confirm", item.start.String() + item.end.String() + item.message, walk.MsgBoxOK)
}

func GetTime(s string) (*time.Time, *time.Time) {
	start := time.Now()
	// 数字だけなら分として扱う
	if _, err := strconv.Atoi(s); err == nil {
		end := start.Add(time.Duration{m * time.Minute})
		return &start, &end
	}
	if d, err := time.ParseDuration(s); err == nil {
		end := start.Add(d)
		return &start, &end
	}
	re := regexp.MustCompile("([0-9]+):([0-9]+)")
	if result := re.MatchString(s); if result {
		hh := time.Duration(re.SubexpNames()[1].Atoi())
		mm := time.Duration(re.SubexpNames()[1].Atoi())
		end := Parse(start.Year(), start.Month(), start.Day, hh, mm, start.Second(), start.NanoSecond(), start.Location())
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
