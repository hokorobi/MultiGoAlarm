package main

import (
	"context"
	_ "embed"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

//go:embed icon/alarm-note.png
var imgAlarmNote []byte

func alarm(s string) {
	var message string

	aw := newAw()

	if s == "" {
		message = "It's Time!"
	} else {
		message = s
	}

	// FIXME: Make a clear icon

	icon, err := walk.NewIconFromImageForDPI(getIcon(imgAlarmNote), 96)
	if err != nil {
		logg(err)
	}

	winsize := declarative.Size{Width: 300, Height: 300}
	// FIXME: too big button
	err = declarative.MainWindow{
		AssignTo: &aw.mw,
		Title:    "Alarm",
		MinSize:  winsize,
		MaxSize:  winsize,
		Size:     winsize,
		Icon:     icon,
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.LinkLabel{
				Text:    message,
				Font:    declarative.Font{Family: "Meiryo", PointSize: 18},
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
		logf(err)
	}

	// Windowスタイルの動的変更　その3 トップレベル表示: Xo式　実験室（labo.xo-ox.net）
	// http://labo.xo-ox.net/article/99823284.html
	//   "生成時にはGWL_EXSTYLEに8(WS_EX_TOPMOST)を加えてやれば良いのだが
	//   一旦生成したWindowに対してsetwindowlongで変更を加えても反映されない｡
	//    setwindowposで-1(HWND_TOPMOST)と-2(HWND_NOTOPMOSTを指定してやる必要がある"
	// "ウインドウサイズ" http://eternalwindows.jp/winbase/window/window13.html
	win.SetWindowPos(aw.mw.Handle(), win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_FRAMECHANGED|win.SWP_NOMOVE|win.SWP_NOSIZE)

	// FIXME: ウィンドウ表示時にフォーカスを移したい
	aw.alarm()

	aw.mw.Run()
}

type alarmWindow struct {
	mw *walk.MainWindow
}

func (aw *alarmWindow) alarm() {
	orgX, orgY := aw.getWindowPos()
	// "Golangで周期的に実行するときのパターン - Qiita" https://qiita.com/tetsu_koba/items/1599408f537cb513b250
	ctx, cancel := context.WithCancel(context.Background())
	t := time.NewTicker(100 * time.Millisecond)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				aw.cycleMove(orgX, orgY)
			}
		}
	}(ctx)
	go func() {
		time.Sleep(2 * time.Second)
		t.Stop()
		cancel()
	}()
}
func (aw *alarmWindow) cycleMove(orgX int32, orgY int32) {
	curX, curY := aw.getWindowPos()
	if orgX == curX && orgY == curY {
		// 1
		aw.moveWindow(orgX+50, orgY-50)
	} else if orgX+50 == curX && orgY-50 == curY {
		// 2
		aw.moveWindow(orgX+100, orgY)
	} else if orgX+100 == curX && orgY == curY {
		// 3
		aw.moveWindow(orgX+50, orgY+50)
	} else {
		aw.moveWindow(orgX, orgY)
	}
}
func (aw *alarmWindow) moveWindow(x int32, y int32) {
	win.SetWindowPos(aw.mw.Handle(), win.HWND_TOPMOST, x, y, 0, 0, win.SWP_FRAMECHANGED|win.SWP_NOSIZE)
}
func (aw *alarmWindow) getWindowPos() (int32, int32) {
	// https://github.com/lxn/walk/blob/55ccb3a9f5c1dae7b1c94f70ea4f9db6afcb5021/form.go#L598-L615
	var r win.RECT
	win.GetWindowRect(aw.mw.Handle(), &r)
	return r.Left, r.Top
}
func newAw() alarmWindow {
	var aw alarmWindow
	var err error

	aw.mw, err = walk.NewMainWindow()
	if err != nil {
		logf(err)
	}

	return aw
}
