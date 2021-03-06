package webrtcSocket

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var conn net.Conn

func Init() {

	var (
		host    = "127.0.0.1"
		port    = "8888"
		remote  = host + ":" + port
		message = "GET /sign_in?chenhui.ma HTTP/1.0\r\n\r\n"
	)

	con, error := net.Dial("tcp", remote)
	conn = con

	if error != nil {

		fmt.Printf("Host not found: %s\n", error)
		os.Exit(1)
	}
	// defer con.Close();
	in, error := con.Write([]byte(message))
	if error != nil {
		fmt.Printf("Error sending data: %s, in: %d\n", error, in)
		os.Exit(2)
	}

	fmt.Println("Connection OK")
	//fmt.Fprintf(con, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {

		fmt.Println(err)
	}
	fmt.Println(status)

}

func SendBuffer(buf []byte) {
	in, error := conn.Write(buf)
	if error != nil {
		fmt.Printf("Error sending data: %s, in: %d\n", error, in)
		os.Exit(2)
	}

	fmt.Println("Connection OK")
	//fmt.Fprintf(con, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {

		fmt.Println(err)
	}
	fmt.Println(status)

}
