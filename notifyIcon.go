package main

import (
	"os"

	"github.com/lxn/walk"
)

func NotifyIcon(app *app) *walk.NotifyIcon {
	// load icon
	icon, err := walk.NewIconFromFile("alarm-check.ico")
	if err != nil {
		Logf(err)
	}

	ni, err := walk.NewNotifyIcon(app.mw)
	if err != nil {
		Logf(err)
	}

	// Set the icon and a tool tip text.
	err = ni.SetIcon(icon)
	if err != nil {
		Logf(err)
	}

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	err = exitAction.SetText("E&xit")
	if err != nil {
		Logf(err)
	}
	exitAction.Triggered().Attach(
		func() {
			ni.Dispose()
			// TODO: Improve exit
			os.Exit(0)
		},
	)
	err = ni.ContextMenu().Actions().Add(exitAction)
	if err != nil {
		Logf(err)
	}

	// FIXME: Open を選択すると再度コンテキストメニューが表示される
	openListWindow := walk.NewAction()
	err = openListWindow.SetText("O&pen")
	if err != nil {
		Logf(err)
	}
	openListWindow.Triggered().Attach(
		func() {
			go ListWindow(*app)
		},
	)
	err = ni.ContextMenu().Actions().Add(openListWindow)
	if err != nil {
		Logf(err)
	}

	// The notify icon is hidden initially, so we have to make it visible.
	err = ni.SetVisible(true)
	if err != nil {
		Logf(err)
	}

	return ni
}
