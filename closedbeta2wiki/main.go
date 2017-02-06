package main

import (
	"fmt"

	config "github.com/agareev/MoexLib/other"
	confluence "github.com/seppestas/go-confluence"
)

func main() {
	var configuration config.Config
	configuration = config.ReadConfig("confg.json")
	z := confluence.BasicAuth("myname", "mypass")
	x, _ := confluence.NewWiki("http://wiki.myurl.com", z)
	ff := make([]string, 3)
	ff[0] = "sss"
	f, _ := x.GetContent("47646880", ff)
	fmt.Println(f.Body)
	fmt.Println(configuration.Server)
}
