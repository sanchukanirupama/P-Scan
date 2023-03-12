package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
    "time"
)

func main() {
    fmt.Print("***************************************************************************** \n")
    fmt.Print("***************************************************************************** \n")
    fmt.Print("                     P-Scan - Concurrent Port Scanner \n")
    fmt.Print("***************************************************************************** \n")
    fmt.Print("***************************************************************************** \n")
    fmt.Print("  \n")
    
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter the address: ")
    address, _ := reader.ReadString('\n')
    address = strings.TrimSpace(address)
    fmt.Printf("Address: %s\n", address)

    fmt.Print("Enter the protocol (tcp/udp): ")
    protocol, _ := reader.ReadString('\n')
    protocol = strings.TrimSpace(protocol)
    fmt.Printf("Protocol: %s\n", protocol)

    fmt.Print("Enter the starting port: ")
    startPortStr, _ := reader.ReadString('\n')
    startPortStr = strings.TrimSpace(startPortStr)
    startPort, err := strconv.Atoi(startPortStr)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Starting port: %d\n", startPort)

    fmt.Print("Enter the ending port: ")
    endPortStr, _ := reader.ReadString('\n')
    endPortStr = strings.TrimSpace(endPortStr)
    endPort, err := strconv.Atoi(endPortStr)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Ending port: %d\n", endPort)

    openPorts := scanPorts(address, protocol, startPort, endPort)
    if len(openPorts) == 0 {
        fmt.Println("No open ports found")
    } else {
        fmt.Printf("Open ports: %v\n", openPorts)
    }

}

// Scan the ports in given range and outputs the open ports as an array
func scanPorts(address string, protocol string, startPort int, endPort int) []int {
    openPorts := []int{}
    results := make(chan int)
    startTime := time.Now()

    for port := startPort; port <= endPort; port++ {
        go func(port int) {
            target := fmt.Sprintf("%s:%d", address, port)
            conn, err := net.Dial(protocol, target)
            if err == nil {
                conn.Close()
                results <- port
            } else {
                results <- 0
            }
        }(port)
    }

    for i := 0; i < endPort-startPort+1; i++ {
        port := <-results
        if port != 0 {
            openPorts = append(openPorts, port)
        }
    }

    endTime := time.Now()
    duration := endTime.Sub(startTime)
    fmt.Print("  \n")
    fmt.Printf("Scanned %d ports in %s\n", endPort-startPort+1, duration)

    return openPorts
}
