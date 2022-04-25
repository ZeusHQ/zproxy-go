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

	var addr = flag.String("addr", "dev.z:80", "The address of the application.")
	var dir = flag.String("dir", path, "The Turborepo project directory to proxy. ")

	flag.Parse()

	fmt.Println("ZProxy Started", *addr)

	handler := AddHosts(*dir)

	OpenHosts(handler)

	defer RemoveHosts(handler)

	log.Fatal(http.ListenAndServe(*addr, handler))
}
