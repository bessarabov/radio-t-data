package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type Channel struct {
	Items []Item `xml:"channel>item"`
}

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type Item struct {
	Description string    `xml:"description"`
	Enclosure   Enclosure `xml:"enclosure"`
	Link        string    `xml:"link"`
	Title       string    `xml:"title"`
}

func parse_mp3(url string) string {

	fmt.Println(url)

	re := regexp.MustCompile(`^.*?([0-9]+)\.mp3$`)

	episode_number := re.FindStringSubmatch(url)[1]

	type RadioTFile struct {
		Number int64  `json:"number"`
		Url    string `json:"url"`
	}

	type RadioTEpisode struct {
		File RadioTFile `json:"file"`
	}

	num, e := strconv.ParseInt(episode_number, 10, 64)
	if e != nil {
		os.Exit(1)
	}

	episode := RadioTEpisode{
		File: RadioTFile{
			Number: num,
			Url:    url,
		},
	}

	var jsonData []byte
	jsonData, err := json.MarshalIndent(episode, "", "    ")
	if err != nil {
		os.Exit(1)
	}

	f, err := os.Create("/data/episodes/" + episode_number + ".json")

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	f.WriteString(string(jsonData) + "\n")

	f.Close()

	return "ok"
}

func main() {

	rss_url := "https://radio-t.com/podcast.rss"

	parse_mp3("http://cdn.radio-t.com/rt_podcast672.mp3")

	resp, err := http.Get(rss_url)

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	channel := Channel{}

	if err := xml.NewDecoder(resp.Body).Decode(&channel); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if len(channel.Items) == 0 {
		fmt.Println("No items")
		os.Exit(1)
	}

	for _, item := range channel.Items {
		parse_mp3(item.Enclosure.Url)
	}
}
