package main

import (
	"fmt"

	"github.com/szabogabriel/foodview/database"
)

func main() {
	fmt.Println("Initializing the application...")
	database.Connect()
}
