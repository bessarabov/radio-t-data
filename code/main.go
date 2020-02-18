package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
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

func main() {

	rss_url := "https://radio-t.com/podcast.rss"

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

	for i, item := range channel.Items {
		fmt.Printf("%v. item title: %v\n", i, item.Enclosure.Url)
	}

}
