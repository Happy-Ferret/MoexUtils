package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"

	moexlib "github.com/agareev/MoexLib/monitoring"
)

func extractCSV(state []byte) map[string]bool {
	//TODO extract haproxy CSV
	a := string(state)
	r := csv.NewReader(strings.NewReader(a))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(record[1], record[17])
	}

	var x map[string]bool
	return x
}

func main() {
	fmt.Println("a")
	m := moexlib.GetAllContents("http://127.0.0.1:9300/haproxy?stats;csv")
	extractCSV(m)
	// extractCSV
	//moexlib.Send2Graphite(delta, key, graphiteHost, graphitePort)
}
