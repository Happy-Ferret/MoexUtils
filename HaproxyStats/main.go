package main

import (
	"fmt"

	moexlib "github.com/agareev/MoexLib/monitoring"
)

func extractCSV(state []byte) map[string]bool {
	//TODO extract haproxy CSV
	var x map[string]bool
	return x
}

func main() {
	fmt.Println("a")
	moexlib.GetAllContents("http://127.0.0.1:9300")
	// extractCSV
	//moexlib.Send2Graphite(delta, key, graphiteHost, graphitePort)
}
