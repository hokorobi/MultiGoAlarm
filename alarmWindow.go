package main

import (
	"log"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func AlarmWindow(s string) {
	var message string

	var mw *walk.MainWindow

	if s == "" {
		message = "It is Time!"
	} else {
		message = s
	}

	if _, err := (declarative.MainWindow{
		AssignTo: &mw,
		Title:    "Alarm",
		MinSize:  declarative.Size{Width: 300, Height: 300},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Label{
				Text: message,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}
