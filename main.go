package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/rodolfoag/gow32"
)

func main() {
	_, err := gow32.CreateMutex("MultiGoAlarm")
	if err != nil {
		// TODO: 引数があったらアラームとして追加
		// fmt.Printf("Error: %d - %s\n", int(err.(syscall.Errno)), err.Error())
		log.Fatal(err)
	}

	app := newApp()

	t := time.NewTicker(time.Second)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				app.update()
			}
		}
	}()

	app.ni = NotifyIcon(app.mw)
	defer app.ni.Dispose()

	ListWindow(app)

	Logg("Run.")
	defer Logg("Stop.")

	if _, err := (declarative.MainWindow{
		AssignTo: &app.mw,
		Title:    "MultiGoAlarm",
		Visible:  false,
	}.Run()); err != nil {
		Logf(err)
	}
}

// App はこのアプリ全体の型
type app struct {
	mw   *walk.MainWindow
	list *AlarmList
	ni   *walk.NotifyIcon
}

func newApp() app {
	var app app
	var err error
	app.mw, err = walk.NewMainWindow()
	if err != nil {
		Logf(err)
	}
	app.list = NewAlarmList()
	return app
}

func (app *app) update() {
	if len(app.list.list) < 1 {
		return
	}

	items := app.list.update()
	app.alarm(items)
}
func (app *app) alarm(items []AlarmItem) {
	for i := range items {
		go AlarmWindow(items[i].Message)
		time.Sleep(100 * time.Millisecond)
	}
}

// https://qiita.com/KemoKemo/items/d135ddc93e6f87008521#comment-7d090bd8afe54df429b9
func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
func getFilename(ext string) string {
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), getFileNameWithoutExt(exec)+ext)
}

func Logg(m interface{}) {
	f, err := os.OpenFile(getFilename(".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("Cannot open log file: " + err.Error())
	}
	defer f.Close()

	log.SetOutput(io.MultiWriter(f, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println(m)
}
func Logf(m interface{}) {
	Logg(m)
	os.Exit(1)
}
