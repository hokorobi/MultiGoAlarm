package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"time"
	"os"
	// "strings"
)

type MyMainWindow struct {
	*walk.MainWindow
	time    *walk.LineEdit
	message *walk.LineEdit
	list    *walk.ListBox
}

type SubWindow struct {
	*walk.MainWindow
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

	for now := range time.Tick(time.Second) {
		fmt.Println(now)
	}

}

func (mw *MyMainWindow) clickAdd() {
	Alarm()
	time := mw.time.Text()
	message := mw.message.Text()
	fmt.Printf("%v %v\n", time, message)
}

func Alarm() {
	sw := &SubWindow{}

	if _, err := (MainWindow{
		AssignTo: &sw.MainWindow,
		Title:    "Alarm",
		MinSize:  Size{300, 300},
		Layout:   VBox{},
		Children: []Widget{
			Label{
				Text: "Name:",
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

func (mw *MyMainWindow) clickQuit() {
	os.Exit(0)
}
