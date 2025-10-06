package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	printerName string
	bind        string
	secret      string
)

func main() {
	printerName = os.Getenv("PRINTER")
	bind = os.Getenv("BIND")
	secret = os.Getenv("SECRET")

	http.HandleFunc("/print-"+secret, handlePrint)

	log.Printf("Label printer API starting on %s", bind)
	log.Printf("Configured printer: %s", printerName)

	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Fatal(err)
	}
}

func handlePrint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create temporary file for the image
	f, err := os.CreateTemp("", "label-*.png")
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	tmpFile := f.Name()
	defer func() {
		f.Close()
		if err := os.Remove(tmpFile); err != nil {
			log.Printf("Failed to remove temporary file %s: %v", tmpFile, err)
		} else {
			log.Printf("Cleaned up temporary file: %s", tmpFile)
		}
	}()

	// Write the request body (image data) to the temp file
	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Printf("Error writing image data: %v", err)
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	// Close file before printing
	if err := f.Close(); err != nil {
		log.Printf("Error closing temp file: %v", err)
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	// Print the label using lp command
	cmd := exec.Command("lp", tmpFile, "-d", printerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error printing label: %v, output: %s", err, string(output))
		http.Error(w, fmt.Sprintf("Failed to print: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully printed label: %s", tmpFile)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Label printed successfully\n"))
}
