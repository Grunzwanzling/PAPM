package main

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"github.com/tobischo/gokeepasslib"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var db *gokeepasslib.Database
var unlocked bool
var cfg Config

func handleSigterm() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			var _ = os.Remove(cfg.Socket)
			println("Terminating")
			os.Exit(1)
		}

	}()

}
func main() {
	handleSigterm()
	cfg = readFlags()
	var err = os.Remove(cfg.Socket)
	l, err := net.ListenUnix("unix", &net.UnixAddr{cfg.Socket, "unix"})
	if err != nil {
		println("Listen error: ", err.Error())
		return
	}
	println("Chmod attempt")
	if err := os.Chmod(cfg.Socket, 0777); err != nil {
		println("Chmod error:" + err.Error())
	} else {
		println("Chmod succes")
	}
	for {
		fd, err := l.AcceptUnix()
		if err != nil {
			println("Accept error: ", err.Error())
			return
		}

		ucred, err := getCredentials(fd)
		if err != nil {

			println("File error: ", err.Error())
		}

		fmt.Printf("peer_pid: %d\n", ucred.Pid)
		proc, err2 := process.NewProcess(ucred.Pid)
		if err2 != nil {

			println("Process error: ", err2.Error())
		}
		//if(ucred.Uid==0){
		println("Execution path:")
		path := getParents(func(proc *process.Process) (string, error) { return proc.Exe() }, proc)
		println("Cmdline:")
		cmd := getParents(func(proc *process.Process) (string, error) { return proc.Cmdline() }, proc)
		go server(fd, path, cmd)
		//}

	}
}

type getParam func(*process.Process) (string, error)

func getParents(f getParam, proc *process.Process) []string {

	var parents []string

	exe, err3 := f(proc)

	if err3 != nil {
		println("Error getting path of process: ", err3.Error())

	} else {
		println("Found: ", exe)
	}
	parents = append(parents, exe)
	for i := 0; i < 10; i++ {

		proc, err4 := proc.Parent()

		if err4 != nil {

			println("Error getting parent: ", err4.Error())

		}
		exe, err3 := f(proc)

		if err3 != nil {
			println("Error getting path of process: ", err3.Error())

		}
		parents = append(parents, exe)
		println("Found: ", exe)

	}
	return parents
}

func getCredentials(conn *net.UnixConn) (*syscall.Ucred, error) {
	f, err := conn.File()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return syscall.GetsockoptUcred(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
}

func server(c net.Conn, exe []string, cmd []string) {
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
			//send(c, "Unlocking\n")

			db, unlockErr = unlockDB(input[1], input[2])
			if unlockErr != nil {

				send(c, "Unlock error: "+unlockErr.Error())
			} else {
				unlocked = true
				send(c, "Unlocked!")
				//	entry := db.Content.Root.Groups[0].Entries[0]
				//	fmt.Println(entry.GetPassword())
			}

		case "lock":
			if !unlocked {

				send(c, "Not unlocked")
				break

			}
			send(c, "Locking\n")
			db.LockProtectedEntries()
			unlocked = false
		case "check":
			if unlocked {

				send(c, "Unlocked\n")
				break

			} else {

				send(c, "Locked\n")
				break
			}
		case "get":
			//			send(c, "Getting")
			if !unlocked {

				send(c, "Not unlocked!")
			}

			root := db.Content.Root
			levels := strings.Split(input[1], "/")
			currentElement := root.Groups[0]
			elem, err := recursiveSearch(currentElement, levels, 0)
			if err == nil {

				found := false
				for _, pair := range elem.Values {

					if pair.Key == "whitelist_path" {
						for _, parent := range exe {
							if contains(strings.Split(pair.Value.Content, ";"), parent) {
								found = true
								println("Found in whitelist")
							}
						}
					}

					if pair.Key == "whitelist_cmd" {
						for _, parent := range cmd {
							if contains(strings.Split(pair.Value.Content, ";"), parent) {
								found = true
								println("Found in whitelist")
							}
						}
					}
				}
				if !found {

					send(c, "Process is not in whitelist")
					break commands
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
