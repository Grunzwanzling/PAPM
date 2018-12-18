package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"github.com/tobischo/gokeepasslib"
	"net"
	"os"
	"strings"
	"syscall"
)

var db *gokeepasslib.Database
var unlocked bool
var socket string

func readFlags() {

	flag.StringVar(&socket, "socket", "~/socket", "a filepath")
	flag.Parse()

	wd, _ := os.Getwd()
	socket = strings.Replace(socket, "~", wd, -1)

}
func main() {
	readFlags()
	l, err := net.ListenUnix("unix", &net.UnixAddr{socket, "unix"})
	if err != nil {
		println("Listen error: ", err.Error())
		return
	}

	for {
		fd, err := l.AcceptUnix()
		if err != nil {
			println("Accept error: ", err.Error())
			return
		}

		ucred, error := getCredentials(fd)
		if error != nil {

			println("File error: ", error.Error())
		}

		fmt.Printf("peer_pid: %d\n", ucred.Pid)
		proc, err2 := process.NewProcess(ucred.Pid)
		if err2 != nil {

			println("Process error: ", err2.Error())
		}

		exe, err3 := proc.Exe()

		if err3 != nil {
			println("Error getting path of process: ", err3.Error())

		}

		println("Execution path: ", exe)

		//if(ucred.Uid==0){
		go server(fd, exe)
		//}

	}
}

func getCredentials(conn *net.UnixConn) (*syscall.Ucred, error) {
	f, err := conn.File()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return syscall.GetsockoptUcred(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
}

func server(c net.Conn, exe string) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		println("Server got: ", string(data))
		input := strings.Split(string(data), ";")

		var unlockErr error

	commands:
		switch command := input[0]; command {

		case "unlock":
			if unlocked {

				send(c, "Already unlocked")
				break

			}
			send(c, "Unlocking\n")

			db, unlockErr = unlockDB(input[1], input[2])
			if unlockErr != nil {

				send(c, "Unlock error: "+unlockErr.Error())
			} else {
				unlocked = true

				entry := db.Content.Root.Groups[0].Entries[0]
				fmt.Println(entry.GetPassword())
			}

		case "lock":
			if !unlocked {

				send(c, "Not unlocked")
				break

			}
			send(c, "Locking\n")
			db.LockProtectedEntries()
			unlocked = false
		case "get":
			//			send(c, "Getting")
			if !unlocked {

				send(c, "Not unlocked!")
				break

			}

			root := db.Content.Root
			levels := strings.Split(input[1], "/")
			currentElement := root.Groups[0]
			elem, err := recursiveSearch(currentElement, levels, 0)
			if err == nil {

				for _, pair := range elem.Values {

					if pair.Key == "whitelist" {

						if !contains(strings.Split(pair.Value.Content, ";"), exe) {

							send(c, "Process is not in whitelist")
							break commands
						}
					}

				}

				send(c, elem.GetPassword())
			} else {
				send(c, err.Error())

			}
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func recursiveSearch(element gokeepasslib.Group, levels []string, lvl int) (gokeepasslib.Entry, error) {

	if lvl+1 == len(levels) {
		for _, elem := range element.Entries {
			println("Searching for entry: " + levels[lvl])
			if elem.GetTitle() == levels[lvl] {
				return elem, nil
			}
		}
	}

	for _, elem := range element.Groups {
		println("Searching for: " + levels[lvl])
		println(elem.Name)
		if elem.Name == levels[lvl] {
			println("Found!")
			return recursiveSearch(elem, levels, lvl+1)

		}
	}
	return gokeepasslib.NewEntry(), errors.New("Not found")

}

func send(c net.Conn, text string) {

	println(text)

	_, err := c.Write([]byte(text))
	if err != nil {
		println("Write error; ", err.Error())
	}

}
func unlockDB(path string, pass string) (*gokeepasslib.Database, error) {

	file, _ := os.Open(path)

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(pass)
	err := gokeepasslib.NewDecoder(file).Decode(db)
	if err != nil {

		return nil, err
	}

	db.UnlockProtectedEntries()
	return db, nil
}
