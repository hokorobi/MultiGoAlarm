package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	// "strings"
)

type MyMainWindow struct {
	*walk.MainWindow
	time    *walk.LineEdit
	message *walk.LineEdit
	list    *walk.ListBox
}

func main() {
	mw := &MyMainWindow{}

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "MultiGoAlarm",
		MinSize:  Size{400, 300},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &mw.time,
					},
					PushButton{
						Text:      "&Add",
						OnClicked: mw.clickAdd,
					},
					PushButton{
						Text:      "&Quit",
						OnClicked: mw.clickQuit,
					},
				},
			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &mw.message,
					},
					ListBox{
						AssignTo: &mw.list,
						Row:      10,
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}

}

func (mw *MyMainWindow) clickAdd() {
	Alarm(mw)
	time := mw.time.Text()
	message := mw.message.Text()
	fmt.Printf("%v %v\n", time, message)
}

func Alarm(owner walk.Form) (int, error) {
	var dlg *walk.Dialog

	return Dialog{
		AssignTo: &dlg,
		Title:    "Alarm",
		MinSize:  Size{300, 300},
		Layout:   VBox{},
		Children: []Widget{
			Label{
				Text: "Name:",
			},
		},
	}.Run(owner)
}

func (mw *MyMainWindow) clickQuit() {
	mw.Close()
}
