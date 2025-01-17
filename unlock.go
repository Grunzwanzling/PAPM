package main

import (
	"bufio"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"io"
	"net"
	"os"
	"time"
)

var cfg Config
var form *tview.Form
var pwField *tview.InputField
var app *tview.Application

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
func readOnce(r io.Reader) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf[:])
	if err != nil {
		println("Read error: ", err.Error())
		os.Exit(1)
	}
	app.Stop()
	fmt.Println(string(buf[0:n]))

}
func sendCommand() {

	c, err := net.Dial("unix", cfg.Socket)
	if err != nil {
		println("Dial error: ", err.Error())
		os.Exit(1)
	}
	defer c.Close()

	go readOnce(c)
	pw := pwField.GetText()
	msg := "unlock;" + cfg.Db + ";" + pw
	_, err = c.Write([]byte(msg))
	if err != nil {
		println("Write error: ", err)
	}
	time.Sleep(1e9)
	app.Stop()
}
func main() {
	cfg = readFlags()
	pwField = tview.NewInputField().
		SetLabel("Password").
		SetMaskCharacter('*')

	app = tview.NewApplication()
	form = tview.NewForm().SetLabelColor(tcell.ColorWhite).
		AddFormItem(pwField).
		AddButton("Unlock", sendCommand).
		AddButton("Quit", func() {
			app.Stop()
		}).
		SetCancelFunc(func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Unlock a database").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}
}
