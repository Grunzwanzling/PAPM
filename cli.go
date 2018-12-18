package main

import (
	"bufio"
	"flag"
	"io"
	"net"
	"os"
	"strings"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		println(string(buf[0:n]))
	}
}

func main() {
	wd, _ := os.Getwd()
	socket := flag.String("socket", "~/socket", "a filepath")
	cmd := flag.String("command", "", "a supported command")
	flag.Parse()
	*socket = strings.Replace(*socket, "~", wd, -1)
	c, err := net.Dial("unix", *socket)
	if err != nil {
		println("Dial error ", err)
	}
	defer c.Close()
	go reader(c)

	if *cmd != "" {

		//	msg := "unlock;/home/max/pass/test.kdbx;test"
		_, err = c.Write([]byte(*cmd))
		if err != nil {
			println("Write error: ", err)
		}

		//			msg = "get;group1/group2/check"
	} else {
		println("Started in CLI mode with unix-socket: " + *socket)
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
