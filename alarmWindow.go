package main

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
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
	// "user interface - How to set window position and make it unresizable in Go walk - Stack Overflow" https://stackoverflow.com/questions/25949966/how-to-set-window-position-and-make-it-unresizable-in-go-walk
	// FIXME: too big button
	// FIXME: too small font
	err := declarative.MainWindow{
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
	}.Create()
	if err != nil {
		Logf(err)
	}

	// Windowスタイルの動的変更　その3 トップレベル表示: Xo式　実験室（labo.xo-ox.net）
	// http://labo.xo-ox.net/article/99823284.html
	//   "生成時にはGWL_EXSTYLEに8(WS_EX_TOPMOST)を加えてやれば良いのだが
	//   一旦生成したWindowに対してsetwindowlongで変更を加えても反映されない｡
	//    setwindowposで-1(HWND_TOPMOST)と-2(HWND_NOTOPMOSTを指定してやる必要がある"
	// "ウインドウサイズ" http://eternalwindows.jp/winbase/window/window13.html
	win.SetWindowPos(mw.Handle(), win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_FRAMECHANGED|win.SWP_NOMOVE|win.SWP_NOSIZE)

	mw.Run()
}
