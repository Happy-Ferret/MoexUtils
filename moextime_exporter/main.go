package main

import (
	"flag"
	"log"
)

var (
	listenPort = ":9262"
)

func init() {
	flag.StringVar(&listenPort, "L", listenPort, "Lister port")
	flag.Parse()
}

func main() {
	log.Println(listenPort)
}
