package main;

import (
		"github.com/shirou/gopsutil/process"
	//	"reflect"
	//	"runtime"
	"net"
	//	"strings"
	//	"bufio"
	"fmt"
	"log"
	//	"os"
	"syscall"
	//	"time"
)

var pipeFile = "pipe.log"

func main() {
	l, err := net.ListenUnix("unix", &net.UnixAddr{"/tmp/echo.sock", "unix"})
	if err != nil {
		println("listen error", err.Error())
		return
	}

	for {
		fd, err := l.AcceptUnix()
		if err != nil {
			println("accept error", err.Error())
			return
		}

		ucred, error := getCredentials(fd)
		if(error != nil){

			println("error2", error.Error())
		}

		fmt.Printf("peer_pid: %d\n", ucred.Pid)
		fmt.Printf("peer_uid: %d\n", ucred.Uid)
		fmt.Printf("peer_gid: %d\n", ucred.Gid)
proc, err2 := process.NewProcess(ucred.Pid)
if(err2 != nil){

println("error3", err2.Error())
}

exe, err3 := proc.Exe()

if(err3 != nil){
	println("error4", err3.Error())


}

println(exe);

//if(ucred.Uid==0){
		go server(fd)
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


