package main

import (
	"github.com/lxn/walk"
)

func NotifyIcon(mw *walk.MainWindow) *walk.NotifyIcon {
	// load icon
	icon, err := walk.NewIconFromFile("MultiGoAlarm.ico")
	if err != nil {
		Logf(err)
	}

	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		Logf(err)
	}

	// Set the icon and a tool tip text.
	if err := ni.SetIcon(icon); err != nil {
		Logf(err)
	}

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		Logf(err)
	}
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		Logf(err)
	}
	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		Logf(err)
	}
	return ni
}
