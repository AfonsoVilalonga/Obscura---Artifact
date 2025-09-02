package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	port       = ":8080"         // Server port
	dataSize   = 1 * 1024 * 1024 // 5 MB of data
	bufferSize = 32 * 1024       // Buffer size for writing chunks
)

func main() {
	http.HandleFunc("/download", dataHandler)

	fmt.Printf("Server is running on http://localhost%s\n", port)
	fmt.Printf("Send a request using: curl http://localhost%s/download\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for download
	w.Header().Set("Content-Disposition", "attachment; filename=\"data.bin\"")
	w.Header().Set("Content-Type", "application/octet-stream")

	fmt.Printf("Sending 5MB of data to %s\n", r.RemoteAddr)

	rand.Seed(time.Now().UnixNano())
	remaining := dataSize
	buffer := make([]byte, bufferSize)

	for remaining > 0 {
		chunkSize := bufferSize
		if remaining < bufferSize {
			chunkSize = remaining
		}

		_, err := w.Write(buffer[:chunkSize])
		if err != nil {
			fmt.Println("Error writing data:", err)
			return
		}

		remaining -= chunkSize
	}

	fmt.Println("Data transfer completed")
}
