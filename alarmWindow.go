package main

import (
	"log"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

// AlarmWindow はアラームウィンドウを表示する関数
func AlarmWindow(s string) {
	var message string

	var mw *walk.MainWindow

	if s == "" {
		message = "It is Time!"
	} else {
		message = s
	}

	// TODO: ウィンドウを動かして目立たせる
	// TODO: テキストの折り返し
	if _, err := (declarative.MainWindow{
		AssignTo: &mw,
		Title:    "Alarm",
		MinSize:  declarative.Size{Width: 300, Height: 300},
		MaxSize:  declarative.Size{Width: 300, Height: 300},
		Size:     declarative.Size{Width: 300, Height: 300},
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.LinkLabel{
				Text: message,
			},
			declarative.PushButton{
				Text:      "&Close",
				OnClicked: func() { mw.Close() },
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}
