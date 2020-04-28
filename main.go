package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-toast/toast"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/rodolfoag/gow32"
)

func main() {
	_, err := gow32.CreateMutex("MultiGoAlarm")
	if err != nil {
		// TODO: 引数があったらアラームとして追加
		// fmt.Printf("Error: %d - %s\n", int(err.(syscall.Errno)), err.Error())
		log.Fatal(err)
	}

	app := newApp()

	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				// Logg("tick")
				app.update()
			}
		}
		// t.Stop()
	}()

	// ni := NotifyIcon(app.mw)
	// defer ni.Dispose()

	// FIXME: Make a clear icon
	icon, err := walk.Resources.Icon("alarm-check.ico")
	if err != nil {
		Logg(err)
	}

	if _, err := (declarative.MainWindow{
		AssignTo: &app.mw,
		Title:    "MultiGoAlarm",
		MinSize:  declarative.Size{Width: 400, Height: 300},
		Size:     declarative.Size{Width: 400, Height: 300},
		Icon:     icon,
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.ListBox{
				AssignTo:        &app.lb,
				Model:           app.list,
				OnItemActivated: app.lbItemActivated,
				Row:             10,
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						Text:      "&Add",
						OnClicked: app.clickAddDlg,
					},
					declarative.PushButton{
						Text:      "&Quit",
						OnClicked: func() { app.mw.Close() },
					},
				},
			},
		},
	}.Run()); err != nil {
		Logf(err)
	}
}

// App はこのアプリ全体の型
type app struct {
	mw   *walk.MainWindow
	lb   *walk.ListBox
	list *AlarmList
}

type additionalAlarmText struct {
	Text string
}

func newApp() app {
	var app app
	var err error
	app.mw, err = walk.NewMainWindow()
	if err != nil {
		Logf(err)
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
func (app *app) clickAddDlg() {
	newText := new(additionalAlarmText)
	cmd, err := additionalDialog(app.mw, newText)
	if err != nil {
		walk.MsgBox(app.mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
		return
	}
	if cmd == walk.DlgCmdOK {
		item := NewAlarmItem(newText.Text)
		if item == nil {
			walk.MsgBox(app.mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
			return
		}
		// debug
		// walk.MsgBox(mw, "confirm", item.start.String()+item.end.String()+item.message, walk.MsgBoxOK)
		app.list.add(*item)

		iconpath, err2 := filepath.Abs("alarm-check.png")
		if err2 != nil {
			Logg(err2)
		}

		notify := toast.Notification{
			AppID:   "MultiGoAlarm",
			Title:   "Add Alarm",
			Icon:    iconpath,
			Message: item.End.Format("15:04") + " " + item.Message,
		}
		err := notify.Push()
		if err != nil {
			Logg(err)
		}

		app.lb.SetModel(app.list)
	}
}
func (app *app) update() {
	if len(app.list.list) < 1 {
		return
	}

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

func getLogfilename() string {
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), filepath.Base(exec)+".log")
}
func Logg(m interface{}) {
	f, err := os.OpenFile(getLogfilename(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("Cannot open log file:" + err.Error())
	}
	defer f.Close()

	log.SetOutput(io.MultiWriter(f, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println(m)
}
func Logf(m interface{}) {
	Logg(m)
	os.Exit(1)
}
