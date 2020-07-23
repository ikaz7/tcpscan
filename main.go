// concurrent network tcp scanner
package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"time"
)

var host = flag.String("host", "localhost", "host to connect to")
var timeout = flag.Duration("t", 10*time.Second, "time to wait before timing out connection.")
var concurrent = flag.Int("c", 100, "number of concurrent connections.")
var e net.Error

func main() {
	flag.Parse()
	// counting semaphore
	var tokens = make(chan struct{}, 200)
	var ch = make(chan struct{})
	start := time.Now()
	for i := 1; i <= 65535; i++ {
		go func(p int) {
			address := fmt.Sprintf("%s:%d", *host, p)
			tokens <- struct{}{}
			conn, err := net.DialTimeout("tcp", address, 10*time.Second)
			if errors.As(err, &e) {
				// hanging connection, likely filtered port
				if e.Timeout() {
					fmt.Printf("scan: %s timed out ...\n", address)
				}
				// closed port
				<-tokens
				ch <- struct{}{}
				return
			}
			fmt.Printf("scan: %s is open\n", address)
			conn.Close()
			<-tokens
			ch <- struct{}{}
		}(i)
	}

	for i := 1; i <= 65535; i++ {
		<-ch
	}

	finish := time.Since(start).Seconds()
	fmt.Printf("tcpscan took %.3fs\n", finish)
}
