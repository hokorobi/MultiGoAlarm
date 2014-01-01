package main

import (
	"log"

	"github.com/lxn/walk"
)

func NotifyIcon(mw *walk.MainWindow) *walk.NotifyIcon {
	// load icon
	icon, err := walk.NewIconFromFile("MultiGoAlarm.ico")
	if err != nil {
		log.Fatal(err)
	}

	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}

	// Set the icon and a tool tip text.
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal(err)
	}
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}
	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}
	return ni
}
