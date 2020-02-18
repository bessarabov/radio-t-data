package main

import (
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
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

func save_file(url string, file_name string) int64 {

	out, err := os.Create(file_name)
	defer out.Close()

	if err != nil {
		os.Exit(1)
	}

	resp, err := http.Get(url)
	if err != nil {
		os.Exit(1)
	}

	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		os.Exit(1)
	}

	return n
}

func parse_mp3(url string) string {

	fmt.Println(url)

	re := regexp.MustCompile(`^.*?([0-9]+)\.mp3$`)

	episode_number := re.FindStringSubmatch(url)[1]

	num, e := strconv.ParseInt(episode_number, 10, 64)
	if e != nil {
		os.Exit(1)
	}

	tmp_file_name := "/tmp/" + episode_number + ".mp3"

	save_file(url, tmp_file_name)

	type RadioTFile struct {
		Number int64  `json:"number"`
		Url    string `json:"url"`
		Md5    string `json:"md5"`
	}

	type RadioTEpisode struct {
		File RadioTFile `json:"file"`
	}

	episode := RadioTEpisode{
		File: RadioTFile{
			Number: num,
			Url:    url,
			Md5:    get_md5_from_file(tmp_file_name),
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

func get_md5_from_file(file_name string) string {
	f, err := os.Open(file_name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func main() {

	// tmp
	os.Exit(1)

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

	for _, item := range channel.Items {
		parse_mp3(item.Enclosure.Url)
	}
}
