package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type pageVersion struct {
	Version struct {
		Number int `json:"number"`
	} `json:"version"`
}

type pageSpace struct {
	Key string `json:"key"`
}

type pageBodyStorage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type pageBody struct {
	Storage pageBodyStorage `json:"storage"`
}

type pageVersionPUT struct {
	Number int `json:"number"`
}

type pageUpdate struct {
	ID      string         `json:"id"`
	Type    string         `json:"type"`
	Tittle  string         `json:"title"`
	Space   pageSpace      `json:"space"`
	Body    pageBody       `json:"body"`
	Version pageVersionPUT `json:"version"`
}

func (t *pageUpdate) JSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func upPageVersion(pageid string) int {
	var output pageVersion
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url+pageid, nil)
	req.SetBasicAuth(login, password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("couldn't connect ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("couldn't read ", err)
	}
	json.Unmarshal(body, &output)
	return (output.Version.Number + 1)
}

func prepare2Wiki(body *bytes.Buffer) []byte {
	version := upPageVersion(pageid)
	vpageSpace := pageSpace{Key: wikiSpaceName}
	vpageBodyStorage := pageBodyStorage{Value: body.String(), Representation: "storage"}
	vpageBody := pageBody{vpageBodyStorage}
	vpageVersionPUT := pageVersionPUT{Number: version}
	data := pageUpdate{
		ID:      pageid,
		Type:    "page",
		Tittle:  wikiPageTittle,
		Space:   vpageSpace,
		Body:    vpageBody,
		Version: vpageVersionPUT,
	}
	output, err := data.JSON()
	if err != nil {
		log.Println("Problem with parse! ", err)
	}
	return output
}

func push2Wiki(body []byte) error {
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", url+pageid, bytes.NewBuffer(body))

	req.SetBasicAuth(login, password)
	req.Header.Set("Content-Type", "application/json")
	// req.Body.Read(body)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
