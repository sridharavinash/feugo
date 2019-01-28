package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/gosimple/slug"
	"github.com/labstack/echo"
)

type SongItem struct {
	Song string `json:"song"`
}

type ResponseData struct {
	Count int        `json:"count"`
	Data  []SongItem `json:"data"`
}

type APIResponse struct {
	ErrorCode    int           `json:"error_code"`
	ErrorMessage string        `json:"error_message"`
	ResponseData *ResponseData `json:"response"`
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Static("/static", "assets")

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e.Renderer = t

	e.GET("/", indexRender)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func indexRender(c echo.Context) error {
	names, err := getNames()
	if err != nil {
		return c.Render(http.StatusInternalServerError, "index", fmt.Sprintf("Error in GET %v", err))
	}
	return c.Render(http.StatusOK, "index.html", names)
}

func getNames() (string, error) {
	apikey := os.Getenv("PHISH_NET_API")
	if apikey == "" {
		fmt.Println("No Phish.net API env variable set")
		os.Exit(1)
	}

	url := fmt.Sprintf("https://api.phish.net/v3/jamcharts/all?apikey=%s", apikey)

	payload := strings.NewReader("{}")

	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", err
	}

	var songs []string
	fullSongs := apiResponse.ResponseData.Data
	for i := 0; i < len(fullSongs); i++ {
		songs = append(songs, makeSlug(fullSongs[i].Song))
	}

	randName := songs[rand.Intn(len(songs))] + "-" + songs[rand.Intn(len(songs))]
	return randName, nil
}

func makeSlug(s string) string {
	tempString := strings.ToLower(s)

	return slug.Make(tempString)
}
