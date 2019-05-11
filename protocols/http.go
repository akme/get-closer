package protocols

import (
	"io"

	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"
)

// HTTPPing measures latency for HTTP request
func HTTPPing(target string) time.Duration {
	//var resp *http.Response
	var body io.Reader
	target = "http://" + target // dirty fix =(
	req, err := http.NewRequest("GET", target, body)
	req.Header.Set(http.CanonicalHeaderKey("User-Agent"), "get-closer")
	if err != nil {
		fmt.Println(err)
		return 0
	}
	// var remoteAddr string
	trace := &httptrace.ClientTrace{
		ConnectStart: func(network, addr string) {
			//	remoteAddr = addr
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	timeout := time.Duration(5 * time.Second)

	startAt := time.Now()
	client := http.Client{Timeout: timeout}
	_, err = client.Do(req)

	endAt := time.Now()

	if err != nil {
		fmt.Println(err)
		return 0
	}
	duration := endAt.UnixNano() - startAt.UnixNano()

	return time.Duration(duration)
}
