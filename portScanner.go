package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// scanPort attempts to establish a TCP connection to the specified address and port.
// If the connection is successful, it sends the port number through the channel and if not it sends a 0
func scanPort(address string, port int, ch chan<- int) {
	// Attempt to establish a TCP connection
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", address, port), time.Second*10)
	if err != nil {
		// If the connection fails, send 0 through the channel
		ch <- 0
		return
	}
	// If the connection succeeds, close it and send the port number through the channel
	conn.Close()
	ch <- port
}

func main() {
	// Check that the correct number of command line arguments were provided
	if len(os.Args) != 4 {
		fmt.Println("Use: ./port-scanner [address] [start port] [end port]")
		os.Exit(1)
	}

	// Parse the command line arguments
	address := os.Args[1]
	start, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid start port")
		os.Exit(1)
	}
	end, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid end port")
		os.Exit(1)
	}

	if end < start {
		fmt.Println("End port must be greater than start port")
		os.Exit(1)
	}

	fmt.Printf("------------Scanning ports---------------- %d-%d on %s...\n", start, end, address)

	// Create a channel for receiving the results of the port scans
	ch := make(chan int)

	// Launch a goroutine for each port to be scanned
	for i := start; i <= end; i++ {
		go scanPort(address, i, ch)
	}

	// Receive the results of the port scans from the channel
	for i := start; i <= end; i++ {
		port := <-ch
		if port != 0 {
			fmt.Printf("Port %d is open\n", port)
		}
	}
}
