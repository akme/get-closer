package protocols

import (
	"fmt"
	"net"
	"time"
)

// TCPPing measures time for openning connection to TCP port
func TCPPing(host string, tcpPort int) time.Duration {
	startTime := time.Now()
	host = fmt.Sprintf("%s:%d", host, tcpPort)
	conn, err := net.Dial("tcp", host)
	endTime := time.Now()
	if err != nil {
		fmt.Println("could not connect to server: ", err)
	}
	defer conn.Close()
	return endTime.Sub(startTime)
}
