package main

import (
	"bytes"
	"flag"
	"html/template"
	"strings"
)

// AllowedClient JSON struct
type AllowedClient struct {
	Name string
	IP   []string
}

var (
	path           = "./file.txt"
	templatePath   = "template.tmpl"
	tempfile       = "state.tmp"
	pageid         = "1111111"
	login          = "login"
	password       = "pasword"
	conflout       bytes.Buffer
	url            = "wikiurl"
	wikiSpaceName  = "WEBDEVOPS"
	wikiPageTittle = "Tittle Description"
)

func (c *AllowedClient) getName() string {
	return c.Name
}

func (c *AllowedClient) getIP() []string {
	return c.IP
}

func makeAllowedClient(body string) AllowedClient {
	mmm := strings.Split(body, "\n")
	name := mmm[0]
	IPs := parseIP(mmm[1:])
	c := AllowedClient{name, IPs}
	return c
}

func main() {
	flag.StringVar(&path, "f", path, "please use the file")
	flag.StringVar(&login, "l", login, "login")
	flag.StringVar(&password, "p", password, "password")
	flag.Parse()

	t, _ := template.ParseFiles(templatePath)
	data := make([]AllowedClient, 1)

	for _, x := range splitFile(readFile(path)) {
		data = append(data, makeAllowedClient(x))
	}
	t.Execute(&conflout, data)

	Joutput := prepare2Wiki(&conflout)
	_ = push2Wiki(Joutput)
}
