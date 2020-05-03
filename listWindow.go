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

func showListWindow(parent app) {

	lw := newListWindow(parent.list)

	// "Golangで周期的に実行するときのパターン - Qiita" https://qiita.com/tetsu_koba/items/1599408f537cb513b250
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
				lw.update()
			}
		}
	}(ctx)

	// FIXME: Make a clear icon
	icon, err := walk.Resources.Icon("alarm-check.ico")
	if err != nil {
		logg(err)
	}

	logg("Run.")
	defer logg("Stop.")

	if _, err := (declarative.MainWindow{
		AssignTo: &lw.mw,
		Title:    "MultiGoAlarm",
		MinSize:  declarative.Size{Width: 400, Height: 300},
		Size:     declarative.Size{Width: 400, Height: 300},
		Icon:     icon,
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.ListBox{
				AssignTo:        &lw.lb,
				Model:           lw.list,
				OnItemActivated: lw.lbItemActivated,
				Row:             10,
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						Text:      "&Add",
						OnClicked: lw.clickAddDlg,
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
		logf(err)
	}
}

type lw struct {
	mw    *walk.MainWindow
	lb    *walk.ListBox
	list  *alarmList
	count int
}

type additionalAlarmText struct {
	Text string
}

func newListWindow(list *alarmList) lw {
	var lw lw
	var err error
	lw.mw, err = walk.NewMainWindow()
	if err != nil {
		logf(err)
	}
	lw.list = list
	return lw
}

func (lw *lw) lbItemActivated() {
	if lw.lb.CurrentIndex() < 0 {
		return
	}

	lw.list.del(lw.lb.CurrentIndex())
	lw.lb.SetModel(lw.list)
}
func (lw *lw) clickAddDlg() {
	newText := new(additionalAlarmText)
	cmd, err := additionalDialog(lw.mw, newText)
	if err != nil {
		walk.MsgBox(lw.mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
		return
	}
	if cmd == walk.DlgCmdOK {
		item := newAlarmItem(newText.Text)
		if item == nil {
			walk.MsgBox(lw.mw, "Error", "Enter valid time", walk.MsgBoxOK|walk.MsgBoxIconError)
			return
		}
		// debug
		// walk.MsgBox(mw, "confirm", item.start.String()+item.end.String()+item.message, walk.MsgBoxOK)
		lw.list.add(*item)

		iconpath, err := filepath.Abs("alarm-check.png")
		if err != nil {
			logg(err)
		}

		notify := toast.Notification{
			AppID:   "MultiGoAlarm",
			Title:   "Add Alarm",
			Icon:    iconpath,
			Message: item.End.Format("15:04") + " " + item.Message,
		}
		err = notify.Push()
		if err != nil {
			logg(err)
		}

		lw.lb.SetModel(lw.list)
	}
}
func (lw *lw) update() {
	// No unnecessary update
	if lw.count == 0 && lw.count == len(lw.list.list) {
		return
	}

	idx := lw.lb.CurrentIndex()
	lw.lb.SetModel(lw.list)
	lw.count = len(lw.list.list)
	err := lw.lb.SetCurrentIndex(idx)
	if err != nil {
		logg(err)
	}

}
