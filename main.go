package main

import (
	"fmt"
	"github.com/Hell077/YandexCalc/internal"
	"log"
)

func main() {
	srv := internal.NewServer(":8000")
	fmt.Println("Server is running at :8000")
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
