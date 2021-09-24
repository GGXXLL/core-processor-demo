package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/ggxxll/core-processor-demo/internal/bootstrap"
)

func main() {
	// setup rand seeder
	rand.Seed(time.Now().UnixNano())

	// Core Bootstrap
	root, shutdown := bootstrap.Bootstrap()

	// Shutdown
	defer shutdown()

	// Run
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
