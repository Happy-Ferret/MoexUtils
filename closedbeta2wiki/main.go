package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	confluence "github.com/seppestas/go-confluence"
)

// AllowedClient JSON struct
type AllowedClient struct {
	Name string
	IP   []string
}

var (
	path     = "./file.txt"
	pageid   = "http://url.com"
	login    = "login"
	password = "password"
)

func readFile(path string) string {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
		return "unknown error"
	}
	str := string(bs)
	return str
}

func splitFile(body string) []string {
	bodys := strings.Split(body, "\n\n")
	return bodys
}

func parseFile(bodys []string) string {
	body := strings.Join(bodys, "\n-------------\n")
	return body
}

func (c *AllowedClient) getName() string {
	return c.Name
}

func makeAllowedClient(body []string) AllowedClient {
	name := body[0]
	IPs := body[0:]
	c := AllowedClient{name, IPs}
	return c
}

// data := url.Values{}
// 	data.Set("expand", strings.Join(expand, ","))
// 	contentEndPoint.RawQuery = data.Encode()

func push2wiki(contentID, login, password string) {
	z := confluence.BasicAuth(login, password)
	// pageid - url2wiki
	x, _ := confluence.NewWiki(pageid, z)
	expand := make([]string, 1)
	expand = append(expand, "title")
	expand = append(expand, "body")
	f, _ := x.GetContent(contentID, expand)
	fmt.Println(f)
}

func main() {
	flag.StringVar(&path, "f", path, "please use the file")
	flag.StringVar(&login, "l", login, "login")
	flag.StringVar(&password, "p", password, "password")
	flag.Parse()
	z := make([]string, 2)
	c := AllowedClient{"petya", z}
	fmt.Println(c.getName())
	// x := splitFile(readFile(path))
	// fmt.Println(parseFile(x))
	// push2wiki("num", login, password)
	// fmt.Println(splitFile(readFile(path)))
	// clients := []AllowedClient
	for _, x := range splitFile(readFile(path)) {
		c := makeAllowedClient(x)
		fmt.Println(c.getName())

	}

}
