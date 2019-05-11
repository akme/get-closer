package protocols

import (
	"github.com/sparrc/go-ping"
	"time"
)

// ICMPPing implements basic ping
func ICMPPing(host string) time.Duration {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		panic(err)
	}
	pinger.Count = 1
	//startTime := time.Now()

	pinger.Run() // blocks until finished
	//	endTime := time.Now()

	stats := pinger.Statistics()
	return stats.AvgRtt
}
