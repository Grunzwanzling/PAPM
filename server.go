package main;

import (
	"github.com/shirou/gopsutil"
	"reflect"
	"runtime"
	"net"
	"strings"
	"bufio"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
)

var pipeFile = "pipe.log"

func main() {
	l, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		println("listen error", err.Error())
		return
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			println("accept error", err.Error())
			return
		}

		ptrVal := reflect.ValueOf(fd)
		val2 := reflect.Indirect(ptrVal)

		// which is a net.conn from which we get the 'fd' field
		fdmember := val2.FieldByName("fd")
		val3 := reflect.Indirect(fdmember)

		// which is a netFD from which we get the 'sysfd' field
		netFdPtr := val3.FieldByName("sysfd")
		fmt.Printf("netFDPtr= %v\n", netFdPtr)

		// which is the system socket (type is plateform specific - Int for linux)
		if runtime.GOOS == "linux" {
			cfd := int(netFdPtr.Int())
			ucred, _ := syscall.GetsockoptUcred(cfd, syscall.SOL_SOCKET, syscall.SO_PEERCRED)

			fmt.Printf("peer_pid: %d\n", ucred.Pid)
			fmt.Printf("peer_uid: %d\n", ucred.Uid)
			fmt.Printf("peer_gid: %d\n", ucred.Gid)

			if(ucred.Uid==0){
				go server(fd)
			}

		}
	}
}

func server(c net.Conn) {

for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		println("Server got:", string(data))
		_, err = c.Write(data)
		if err != nil {
			log.Fatal("Writing client error: ", err)
		}
	}

}


