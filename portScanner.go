package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: ./P-SCAN IP_ADDRESS START_PORT END_PORT")
		os.Exit(1)
	}

	ip := os.Args[1]
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
	if start > end {
		fmt.Println("Start port must be less than or equal to end port")
		os.Exit(1)
	}

	// Use a wait group to track the status of the concurrent goroutines
	var wg sync.WaitGroup

	// Scan the specified range of ports
	for port := start; port <= end; port++ {
		wg.Add(1)

		// Launch a goroutine to scan the port
		go func(port int) {
			defer wg.Done() // Decrement the wait group counter when the goroutine finishes

			// Try to connect to the host and port
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
			if err != nil {
				return
			} else {
				// If there is no error, the port is open
				fmt.Printf("%d: open\n", port)

				// Look up the service name for the port
				service, err := net.LookupPort("tcp", strconv.Itoa(port))
				if err != nil {
					fmt.Println("Unable to lookup service name")
				} else {
					fmt.Printf("Service name: %s\n", service)
				}

				// Remember to close the connection when finished
				conn.Close()
			}
		}(port) // Pass the current port number to the goroutine
	}

	// Wait for all goroutines to finish before exiting
	wg.Wait()
}
