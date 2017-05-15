package main

import (
	"io/ioutil"
	"log"
	"strings"
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

func parseIP(body []string) []string {
	var parsedBody []string
	for _, x := range body {
		reg := strings.Replace(x, "allow ", "", 1)
		reg = strings.Replace(reg, ";", "", 1)
		parsedBody = append(parsedBody, reg)
	}
	return parsedBody
}
