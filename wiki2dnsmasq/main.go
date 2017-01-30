package main

import (
	"fmt"

	confluence "github.com/seppestas/go-confluence"
)

func main() {
	fmt.Println("a")
	z := confluence.BasicAuth("myname", "mypass")
	x, _ := confluence.NewWiki("http://wiki.myurl.com", z)
	ff := make([]string, 3)
	ff[0] = "sss"
	f, _ := x.GetContent("47646880", ff)
	fmt.Println(f.Body)
}
