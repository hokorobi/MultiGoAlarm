package main

import (
	"log"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

// AlarmWindow はアラームウィンドウを表示する関数
func additionalDialog(owner walk.Form, s *additionalAlarmText) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	return declarative.Dialog{
		AssignTo:      &dlg,
		Title:         "Add alarm. format: time [message]",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: declarative.DataBinder{
			AssignTo:       &db,
			Name:           "alarm",
			DataSource:     s,
			ErrorPresenter: declarative.ToolTipErrorPresenter{},
		},
		MinSize: declarative.Size{Height: 300, Width: 300},
		Layout:  declarative.VBox{},
		Children: []declarative.Widget{
			declarative.LineEdit{
				Text: declarative.Bind("Text"),
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.HSpacer{},
					declarative.PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}

							dlg.Accept()
						},
					},
					declarative.PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
