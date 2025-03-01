package gotcpscanner

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func worker(ports chan int, res chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range ports {
		address := "localhost:" + fmt.Sprintf("%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Printf("Failed to connect to port %d: %v\n", p, err)
			res <- 0
			continue
		}
		conn.Close()
		res <- p
	}
}

const maxWorkers = 100

func Scanning(startPort, endPort int) []int {
	start := time.Now()

	ports := make(chan int, 100)
	results := make(chan int, 1024)
	openports := make([]int, 0)

	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(ports, results, &wg)
	}

	go func() {
		for i := startPort; i < endPort; i++ {
			ports <- i
		}
		close(ports)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 1; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	sort.Ints(openports)

	elapsed := time.Since(start)
	fmt.Printf("Took %s\n", elapsed)

	return openports
}
