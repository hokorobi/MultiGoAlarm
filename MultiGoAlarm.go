package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"fmt"
	"log"
	"strings"
)

type MyMainWindow struct {
	*walk.MainWindow
	searchBox *walk.LineEdit
	textArea *walk.TextEdit
	results *walk.ListBox
}

func main() {
	mw := &MyMainWindow {}
	
	if _, err := (MainWindow {
		AssignTo: &mw.MainWindow,
		Title: "MultiGoAlarm",
		MinSize: Size {300, 400},
		Layout: VBox {},
		Children: []Widget {
			GroupBox {
				Layout: HBox {},
				Children: []Widget {
					LineEdit {
						AssignTo: &mw.searchBox,
					},
					PushButton {
						Text: "検索",
						OnClicked: mw.clicked,
					},
				},
			},
			TextEdit {
				AssignTo: &mw.textArea,
			},
			ListBox {
				AssignTo: &mw.results,
				Row: 5,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}

}


