package main

import (
	"log"

	"example/internal/app"
)

func main() {
	log.Fatal(app.New().Listen(":3000"))
}
