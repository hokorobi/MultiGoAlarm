package main

import (
	"log"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type SubWindow struct {
	*walk.MainWindow
}

func AlarmWindow(s string) {
	var message string

	sw := &SubWindow{}

	lock.RLock()
	defer lock.RUnlock()

	if s == "" {
		message = "It is Time!"
	} else {
		message = s
	}

	if _, err := (declarative.MainWindow{
		AssignTo: &sw.MainWindow,
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
