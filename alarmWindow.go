package main

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

// AlarmWindow はアラームウィンドウを表示する関数
func AlarmWindow(s string) {
	var message string

	var mw *walk.MainWindow

	if s == "" {
		message = "It's Time!"
	} else {
		message = s
	}

	winsize := declarative.Size{Width: 300, Height: 300}
	// TODO: ウィンドウを動かして目立たせる
	// TODO: 最前面にウィンドウを表示
	// FIXME: too big button
	// FIXME: too small font
	if _, err := (declarative.MainWindow{
		AssignTo: &mw,
		Title:    "Alarm",
		MinSize:  winsize,
		MaxSize:  winsize,
		Size:     winsize,
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.LinkLabel{
				Text:    message,
				MaxSize: declarative.Size{Width: 300, Height: 0},
			},
			declarative.VSpacer{},
			declarative.PushButton{
				Text:      "&Close",
				OnClicked: func() { mw.Close() },
			},
		},
	}.Run()); err != nil {
		Logf(err)
	}
}
