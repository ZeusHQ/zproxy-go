package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	var addr = flag.String("addr", "localhost:80", "The address of the application.")
	var dir = flag.String("dir", path, "The Turborepo project directory to proxy. ")

	flag.Parse()

	fmt.Println("ZProxy Started", *addr)

	handler := CreateHandler()
	handler.AddMonorepoHosts(*dir)

	OpenHosts(handler)

	defer RemoveHosts(handler)

	log.Fatal(http.ListenAndServe(*addr, handler))
}
