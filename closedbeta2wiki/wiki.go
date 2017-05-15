package main

import (
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
	return output.Version.Number + 1
}

func push2Wiki(body string) error {
	type Updater struct {
		id     int
		Type   string `json:"type"`
		tittle string
		space  struct {
			key string
		}
		Body struct {
			storage struct {
				value          string
				representation string
			}
		} `json:"body"`
		Version struct {
			number int
		} `json:"version"`
	}
	version := upPageVersion(pageid)
	data := Updater{
		id:     pageid,
		Type:   "page",
		tittle: "new page",
		space{
			key: "WEBDEVOPS",
		},
		Body{
			storage{
				value:          body,
				representation: "storage",
			},
		},
		Version{
			number: version,
		},
	}
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", url+pageid, nil)
	req.SetBasicAuth(login, password)
	req.Header.Set("Content-Type", "application/json")
	_, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(body, x)
	return nil
}

// data := url.Values{}
// 	data.Set("expand", strings.Join(expand, ","))
// 	contentEndPoint.RawQuery = data.Encode()

// func push2wiki(contentID, login, password string) {
// 	z := confluence.BasicAuth(login, password)
// 	// pageid - url2wiki
// 	x, _ := confluence.NewWiki(pageid, z)
// 	expand := make([]string, 1)
// 	expand = append(expand, "title")
// 	expand = append(expand, "body")
// 	f, _ := x.GetContent(contentID, expand)
// 	fmt.Println(f)
// }
