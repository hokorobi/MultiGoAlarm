package main

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-toast/toast"
	"github.com/hokorobi/go-utils/logutil"
	"github.com/lxn/win"
	"github.com/rodolfoag/gow32"
)

type app struct {
	list *alarmList
}

func main() {
	_, err := gow32.CreateMutex("MultiGoAlarm")
	if err != nil {
		if len(os.Args) == 0 {
			os.Exit(1)
		}

		item := newAlarmItem(strings.Join(os.Args[1:], " "))
		if item == nil {
			win.MessageBox(
				win.HWND(0),
				UTF16PtrFromString("Error: Enter valid time format:"+strings.Join(os.Args[1:], " ")),
				UTF16PtrFromString("URL を開く周期を指定してください"),
				win.MB_OK+win.MB_ICONEXCLAMATION)
			os.Exit(1)
		}
		templist := newAlarmList()
		templist.add(*item)
		notification(*item)
		os.Exit(0)
	}

	logutil.PrintTee("Start")
	defer logutil.PrintTee("End")

	app := newApp()

	sc := gocron.NewScheduler(time.UTC)
	sc.Every(5).Seconds().Do(app.update)
	sc.StartBlocking()
}

func newApp() app {
	var app app
	app.list = newAlarmList()
	return app
}

func (app *app) update() {
	items := app.list.update()
	app.alarm(items)

}
func (app *app) alarm(items []alarmItem) {
	for i := range items {
		go alarm(items[i].Message)
		logg("Alarm: " + items[i].End.Format("15:04:05") + " " + items[i].Message)
		time.Sleep(100 * time.Millisecond)
	}
}

// "Go から Windows の MessageBox を呼び出す - Qiita" https://qiita.com/manymanyuni/items/867d7e0112ce22dec6d5
func UTF16PtrFromString(s string) *uint16 {
	result, _ := syscall.UTF16PtrFromString(s)
	return result
}

func execSchedule() {
	logutil.PrintTee("yahoo")
}

func parseInTokyo(layout string, value string) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	t, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return t, err
	}
	return t, nil
}

func notification(item alarmItem) {
	notify := toast.Notification{
		AppID:   "MultiGoAlarm",
		Title:   "Add Alarm",
		Message: item.End.Format("15:04:05") + " " + item.Message,
	}
	err := notify.Push()
	if err != nil {
		logg(err)
	}
}

func getIcon(icon []byte) image.Image {
	img, err := png.Decode(bytes.NewReader(icon))
	if err != nil {
		logf(err)
	}
	return img
}

// https://qiita.com/KemoKemo/items/d135ddc93e6f87008521#comment-7d090bd8afe54df429b9
func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
func getFilename(ext string) string {
	exec, _ := os.Executable()
	return filepath.Join(filepath.Dir(exec), getFileNameWithoutExt(exec)+ext)
}

func logg(m interface{}) {
	f, err := os.OpenFile(getFilename(".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("Cannot open log file: " + err.Error())
	}
	defer f.Close()

	log.SetOutput(io.MultiWriter(f, os.Stderr))
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println(m)
}
func logf(m interface{}) {
	logg(m)
	os.Exit(1)
}
