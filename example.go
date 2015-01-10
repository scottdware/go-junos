package main

import (
	"fmt"
	"github.com/scottdware/go-junos"
	"log"
	"os"
)

var (
	host     = os.Args[1]
	user     = os.Args[2]
	password = os.Args[3]
)

func main() {
	// Establish our session
	jnpr := junos.NewSession(host, user, password)

	// Lock config
	err := jnpr.Lock()
	if err != nil {
		log.Fatal(err)
	}

	// Unlock config
	err = jnpr.Unlock()
	if err != nil {
		log.Fatal(err)
	}

	// Rollback diff compare
	diff, err := jnpr.RollbackDiff(3)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(diff)

	// View rescue config
	rescue, err := jnpr.GetRescueConfig()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(rescue)

	// Show command
	output, err := jnpr.Command("show version", "text")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(output)

	// Close the connection
	jnpr.Close()
}
