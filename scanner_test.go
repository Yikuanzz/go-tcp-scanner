package gotcpscanner

import (
	"testing"
	"time"
)

// TestScanning tests the Scanning function with a known range of ports.
func TestScanning(t *testing.T) {
	// Define the port range to scan
	startPort := 1
	endPort := 100 // Limiting to 100 for faster testing

	// Call the Scanning function
	openPorts := Scanning(startPort, endPort)

	// Check if the result is a slice of integers
	if len(openPorts) == 0 {
		t.Errorf("Expected openPorts to be a non-empty slice, got %v", openPorts)
	}

	// Print the open ports for debugging
	for _, port := range openPorts {
		t.Logf("Port %d is open", port)
	}

	// Example: Check if a specific port is in the result (if known)
	// For example, port 80 is often open on public servers like scanme.nmap.org
	expectedPort := 80
	found := false
	for _, port := range openPorts {
		if port == expectedPort {
			found = true
			break
		}
	}

	if !found {
		t.Logf("Port %d was not found in the open ports list. This might be expected if the server configuration changes.", expectedPort)
	}

	// Additional checks can be added as needed
}

// TestScanningTimeout tests the Scanning function to ensure it completes within a reasonable time.
func TestScanningTimeout(t *testing.T) {
	// Define the port range to scan
	startPort := 1
	endPort := 100 // Limiting to 100 for faster testing

	// Create a channel to signal completion
	done := make(chan bool)

	// Start the Scanning function in a goroutine
	go func() {
		Scanning(startPort, endPort)
		done <- true
	}()

	// Set a timeout
	select {
	case <-done:
		// Scanning completed successfully
	case <-time.After(5 * time.Second):
		// Scanning took too long
		t.Errorf("Scanning did not complete within 5 seconds")
	}
}
