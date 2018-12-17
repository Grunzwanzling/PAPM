package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		println("Client got:", string(buf[0:n]))
	}
}

func main() {
	c, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		log.Fatal("Dial error", err)
	}
	defer c.Close()
	for {
		go reader(c)
		//		msg := "unlock;/home/max/pass/test.kdbx;test"
		//		_, err = c.Write([]byte(msg))
		//		if err != nil {
		//			log.Fatal("Write error:", err)
		//		}
		//		println("Client sent:", msg)
		//		time.Sleep(1e9)
		//
		//		msg = "get;group1/group2/check"
		//		_, err = c.Write([]byte(msg))
		//		if err != nil {
		//			log.Fatal("Write error:", err)
		//		}
		//		println("Client sent:", msg)
		//		time.Sleep(1e9)

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			var text = scanner.Text()
			fmt.Println(text)
			_, err = c.Write([]byte(text))
			if err != nil {
				log.Fatal("Write error:", err)
			}

		}

		if scanner.Err() != nil {
			log.Fatal("Scanner error:", scanner.Err())
		}

	}
}
