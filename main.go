package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	_, err := net.Dial("tcp", "scanme.nmap.org:80")
	if err != nil {
		fmt.Fprintf(os.Stderr, "scan: %v\n", err)
	}

	fmt.Println("Connection successful.")
}
