package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/go-toast/toast"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func ListWindow(parent app) {

	app := newListWindow(parent.list)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t := time.NewTicker(time.Second)
	defer t.Stop()
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				// Logg("Stop go func().")
				return
			case <-t.C:
				// Logg("tick")
				app.update()
			}
		}
	}(ctx)

	// FIXME: Make a clear icon
	icon, err := walk.Resources.Icon("alarm-check.ico")
	if err != nil {
		Logg(err)
	}

	Logg("Run.")
	defer Logg("Stop.")

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
				Model:           parent.list,
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
						Text: "&Quit",
						OnClicked: func() {
							parent.ni.Dispose()
							// TODO: Improve exit
							// parent.mw.Close() では終了できないみたい
							os.Exit(0)
						},
					},
				},
			},
		},
	}.Run()); err != nil {
		Logf(err)
	}
}

type lw struct {
	mw   *walk.MainWindow
	lb   *walk.ListBox
	list *AlarmList
}

type additionalAlarmText struct {
	Text string
}

func newListWindow(list *AlarmList) lw {
	var app lw
	var err error
	app.mw, err = walk.NewMainWindow()
	if err != nil {
		Logf(err)
	}
	app.list = list
	return app
}

func (app *lw) lbItemActivated() {
	if app.lb.CurrentIndex() < 0 {
		return
	}

	app.list.del(app.lb.CurrentIndex())
	app.lb.SetModel(app.list)
}
func (app *lw) clickAddDlg() {
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
func (app *lw) update() {
	if len(app.list.list) < 1 {
		return
	}

	idx := app.lb.CurrentIndex()
	app.lb.SetModel(app.list)
	err := app.lb.SetCurrentIndex(idx)
	if err != nil {
		Logg(err)
	}

}
