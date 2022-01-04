package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Testing purposes only
func main() {
	resp, err := http.Get(os.Getenv("BOT_ADDR") + "/getUpdates")
	if err != nil {
		log.Fatalf("unexpected error while getting updates: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Fatalf("unexpected error while reading response: %v", err)
	}

	fmt.Print(string(body))
}
