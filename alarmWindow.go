package main

import (
	"context"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

// AlarmWindow はアラームウィンドウを表示する関数
func AlarmWindow(s string) {
	var message string

	aw := newAw()

	if s == "" {
		message = "It's Time!"
	} else {
		message = s
	}

	winsize := declarative.Size{Width: 300, Height: 300}
	// FIXME: too big button
	// FIXME: too small font
	err := declarative.MainWindow{
		AssignTo: &aw.mw,
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
				OnClicked: func() { aw.mw.Close() },
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
	win.SetWindowPos(aw.mw.Handle(), win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_FRAMECHANGED|win.SWP_NOMOVE|win.SWP_NOSIZE)

	aw.alarm()

	aw.mw.Run()
}

type aw struct {
	mw   *walk.MainWindow
	orgX int32
	orgY int32
}

func (aw *aw) alarm() {
	aw.orgX, aw.orgY = aw.getWindowPos()
	// "Golangで周期的に実行するときのパターン - Qiita" https://qiita.com/tetsu_koba/items/1599408f537cb513b250
	ctx, cancel := context.WithCancel(context.Background())
	t := time.NewTicker(100 * time.Millisecond)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				aw.cycleMove()
			}
		}
	}(ctx)
	go func() {
		time.Sleep(time.Second)
		t.Stop()
		cancel()
	}()
}
func (aw *aw) cycleMove() {
	curX, curY := aw.getWindowPos()
	if aw.orgX == curX {
		if aw.orgY == curY {
			aw.moveWindow(aw.orgX-50, aw.orgY)
		} else if aw.orgY+50 == curY {
			aw.moveWindow(aw.orgX, aw.orgY-50)
		} else if aw.orgY-50 == curY {
			aw.moveWindow(aw.orgX+50, aw.orgY)
		} else {
			aw.moveWindow(aw.orgX, aw.orgY)
		}
	} else if aw.orgX+50 == curX {
		aw.moveWindow(aw.orgX, aw.orgY)
	} else if aw.orgX-50 == curX {
		aw.moveWindow(aw.orgX, aw.orgY-50)
	}
}
func (aw *aw) moveWindow(x int32, y int32) {
	win.SetWindowPos(aw.mw.Handle(), win.HWND_TOPMOST, x, y, 0, 0, win.SWP_FRAMECHANGED|win.SWP_NOSIZE)
}
func (aw *aw) getWindowPos() (int32, int32) {
	// https://github.com/lxn/walk/blob/55ccb3a9f5c1dae7b1c94f70ea4f9db6afcb5021/form.go#L598-L615
	var r win.RECT
	win.GetWindowRect(aw.mw.Handle(), &r)
	return r.Left, r.Top
}
func newAw() aw {
	var aw aw
	var err error

	aw.mw, err = walk.NewMainWindow()
	if err != nil {
		Logf(err)
	}

	return aw
}
