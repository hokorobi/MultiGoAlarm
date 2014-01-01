package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	// "strings"
	"os"
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
		MinSize:  Size{300, 400},
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
			LineEdit{
				AssignTo: &mw.message,
			},
			ListBox{
				AssignTo: &mw.list,
				Row:      10,
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}

}

func (mw *MyMainWindow) clickAdd() {
	time := mw.time.Text()
	message := mw.message.Text()
	fmt.Printf("%v %v\n", time, message)
}

func (mw *MyMainWindow) clickQuit() {
	os.Exit(0)
}
