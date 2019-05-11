package protocols

import (
	"fmt"
	"net"
	"time"
)

// TCPPing measures time for openning connection to TCP port
func TCPPing(host string) time.Duration {
	startTime := time.Now()
	conn, err := net.Dial("tcp", host+":53")
	endTime := time.Now()
	if err != nil {
		fmt.Println("could not connect to server: ", err)
	}
	defer conn.Close()
	return endTime.Sub(startTime)
}
