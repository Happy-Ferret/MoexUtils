package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func setMtime(updateTime time.Time) {
	mtime := time.Time(updateTime)
	atime := time.Time(updateTime)
	if err := os.Chtimes(tempfile, atime, mtime); err != nil {
		log.Fatal(err)
	}
}

func checkModify() bool {
	fileState1, err := os.Stat(path)
	if err != nil {
		return true
	}
	fileState2, err := os.Stat(tempfile)
	if err != nil {
		setMtime(fileState1.ModTime())
		return true
	}
	if fileState1.ModTime().Format("20060102150405") != fileState2.ModTime().Format("20060102150405") {
		setMtime(fileState1.ModTime())
		return true
	}
	return false
}

func readFile(path string) string {
	if checkModify() == false {
		log.Println("sended")
		os.Exit(0)
	}
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
