package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type SetListItem struct {
	SetlistDataHTML string `json:"setlistdata"`
}

type ResponseData struct {
	Count int           `json:"count"`
	Data  []SetListItem `json:"data"`
}

type APIResponse struct {
	ErrorCode    int           `json:"error_code"`
	ErrorMessage string        `json:"error_message"`
	ResponseData *ResponseData `json:"response"`
}

func main() {
	apikey := os.Getenv("PHISH_NET_API")

	url := fmt.Sprintf("https://api.phish.net/v3/setlist/random?apikey=%s", apikey)

	payload := strings.NewReader("{}")

	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		fmt.Println("Error creating request")
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error getting response")
		os.Exit(1)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Printf("Error Unmarshalling %v\n", err)
	}

	doc, err := html.Parse(strings.NewReader(apiResponse.ResponseData.Data[0].SetlistDataHTML))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var f func(*html.Node)
	var songs []string
	pattern := regexp.MustCompile("http://phish.net/song/(.*)")
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				c := pattern.FindStringSubmatch(attr.Val)
				if c != nil {
					songs = append(songs, c[1])
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	randName := songs[rand.Intn(len(songs))] + "-" + songs[rand.Intn(len(songs))]
	fmt.Println(randName)

}
