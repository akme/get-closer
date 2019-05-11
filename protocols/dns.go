package protocols

import (
	"context"
	"fmt"
	"net"
	"time"
)

// UseCustomDNS set custom DNS resolver
func UseCustomDNS(dns []string) {

	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (conn net.Conn, err error) {
			for _, addr := range dns {
				if conn, err = net.Dial("udp", addr+":53"); err != nil {
					continue
				} else {
					return conn, nil
				}
			}
			return
		},
	}
	net.DefaultResolver = &resolver
}

//DNSPing set timeout for DNS request
func DNSPing(host string) time.Duration {
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel() // important to avoid a resource leak
	r := net.DefaultResolver
	startTime := time.Now()
	_, err := r.LookupHost(ctx, host)
	endTime := time.Now()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(addrs)
	return endTime.Sub(startTime)
}
