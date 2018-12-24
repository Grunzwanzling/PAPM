package main

import (
	"bufio"
	"flag"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os"
	"strings"
)

var socket string
var database_list string

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
func readFlags() {

	flag.StringVar(&socket, "socket", "./socket", "a filepath")
	flag.StringVar(&database_list, "db_list", "./db_list", "a filepath")
	flag.Parse()

	wd, _ := os.Getwd()
	socket = strings.Replace(socket, ".", wd, -1)
	database_list = strings.Replace(database_list, ".", wd, -1)

}
func main() {
	readFlags()
	dbs, err := readLines(database_list)
	if err != nil {
		println("Read error: ", err.Error())
		os.Exit(0)
	}

	app := tview.NewApplication()
	form := tview.NewForm().SetLabelColor(tcell.ColorWhite).
		AddDropDown("Database", dbs, 0, nil).
		AddPasswordField("Password", "", 50, '*', nil).
		AddButton("Unlock", nil).
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
