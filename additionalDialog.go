package main

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func additionalDialog(owner walk.Form, s *additionalAlarmText) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	// TODO: smaller button
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
		MinSize: declarative.Size{Height: 50, Width: 300},
		Layout:  declarative.VBox{},
		Children: []declarative.Widget{
			declarative.LineEdit{
				Text: declarative.Bind("Text"),
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							err := db.Submit()
							if err != nil {
								logg(err)
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
