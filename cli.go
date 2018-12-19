package main

import (
	"bufio"
	"flag"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var socket string
var cmd string

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			println("Read error: ", err.Error())
			os.Exit(1)
		}
		println(string(buf[0:n]))
	}
}

func readOnce(r io.Reader) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf[:])
	if err != nil {
		println("Read error: ", err.Error())
		os.Exit(1)
	}
	println(string(buf[0:n]))
}

func readFlags() {

	flag.StringVar(&socket, "socket", "~/socket", "a filepath")
	flag.StringVar(&cmd, "command", "", "a supported command")
	flag.Parse()

	wd, _ := os.Getwd()
	socket = strings.Replace(socket, "~", wd, -1)

}
func main() {
	readFlags()
	println(os.Getpid())
	c, err := net.Dial("unix", socket)
	if err != nil {
		println("Dial error: ", err.Error())
		os.Exit(1)
	}
	defer c.Close()

	if cmd != "" {
		go readOnce(c)
		//	msg := "unlock;/home/max/pass/test.kdbx;test"
		_, err = c.Write([]byte(cmd))
		if err != nil {
			println("Write error: ", err)
		}
		time.Sleep(1e9)
		os.Exit(1)
		//			msg = "get;group1/group2/check"
	} else {

		go reader(c)
		println("Started in CLI mode with unix-socket: " + socket)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			var text = scanner.Text()
			_, err = c.Write([]byte(text))
			if err != nil {
				println("Write error: ", err)
			}

		}
		if scanner.Err() != nil {
			println("Scanner error: ", scanner.Err())
		}

	}
}
