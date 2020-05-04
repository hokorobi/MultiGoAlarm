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
	"github.com/lxn/win"
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

	logg("Run.")
	defer logg("Stop.")

	if len(os.Args) > 1 {
		item := newAlarmItem(strings.Join(os.Args[1:], " "))
		if item == nil {
			walk.MsgBox(app.mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
		} else {
			app.list.add(*item)
			notification(*item)
		}
	}

	err = declarative.MainWindow{
		AssignTo: &app.mw,
		Title:    "MultiGoAlarm",
		MinSize:  declarative.Size{Width: 400, Height: 300},
		Size:     declarative.Size{Width: 400, Height: 300},
		Visible:  false,
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.ListBox{
				AssignTo:        &app.lb,
				Model:           app.list,
				OnItemActivated: app.lbItemActivated,
				Row:             10,
			},
			declarative.PushButton{
				Text:      "&Add",
				OnClicked: app.clickAddDlg,
			},
		},
		// Minimize to hide
		OnSizeChanged: func() {
			if win.IsIconic(app.mw.Handle()) {
				app.mw.Hide()
			}
		},
	}.Create()
	if err != nil {
		logf(err)
	}

	// https://github.com/lxn/walk/issues/127
	app.addNotifyIcon()
	app.mw.Run()
}

type app struct {
	mw    *walk.MainWindow
	lb    *walk.ListBox
	list  *alarmList
	ni    *walk.NotifyIcon
	count int
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

	if !app.mw.Visible() {
		return
	}

	if app.count == 0 && app.count == len(app.list.list) {
		return
	}

	idx := app.lb.CurrentIndex()
	app.lb.SetModel(app.list)
	app.count = len(app.list.list)
	err := app.lb.SetCurrentIndex(idx)
	if err != nil {
		logg(err)
	}
}
func (app *app) alarm(items []alarmItem) {
	for i := range items {
		go alarm(items[i].Message)
		logg("Alarm: " + items[i].End.Format("15:04:05") + " " + items[i].Message)
		time.Sleep(100 * time.Millisecond)
	}
}
func (app *app) lbItemActivated() {
	if app.lb.CurrentIndex() < 0 {
		return
	}

	app.list.del(app.lb.CurrentIndex())
	app.lb.SetModel(app.list)
}
func (app *app) clickAddDlg() {
	newText := new(additionalAlarmText)
	cmd, err := additionalDialog(app.mw, newText)
	if err != nil {
		walk.MsgBox(app.mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
		return
	}
	if cmd != walk.DlgCmdOK {
		return
	}

	item := newAlarmItem(newText.Text)
	if item == nil {
		walk.MsgBox(app.mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
		return
	}
	app.list.add(*item)
	notification(*item)
	// app.lb.SetModel は app.update() で反映
}
func (app *app) addNotifyIcon() {
	var err error
	app.ni, err = walk.NewNotifyIcon(app.mw)
	if err != nil {
		logf(err)
	}

	// FIXME: Make a clear icon
	icon, err := walk.Resources.Icon("alarm-check.ico")
	if err != nil {
		logg(err)
	}
	err = app.mw.SetIcon(icon)
	if err != nil {
		logg(err)
	}
	err = app.ni.SetIcon(icon)
	if err != nil {
		logg(err)
	}

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	err = exitAction.SetText("E&xit")
	if err != nil {
		logf(err)
	}
	exitAction.Triggered().Attach(
		func() {
			app.ni.Dispose()
			// TODO: Improve exit
			os.Exit(0)
		},
	)
	err = app.ni.ContextMenu().Actions().Add(exitAction)
	if err != nil {
		logf(err)
	}

	app.ni.SetVisible(true)
	app.ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton {
			app.mw.Show()
			win.ShowWindow(app.mw.Handle(), win.SW_RESTORE)
		}
	})

}

type additionalAlarmText struct {
	Text string
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
