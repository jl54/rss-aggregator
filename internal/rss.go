package internal

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type RssFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RssItem `xml:"item"`
	} `xml:"channel"`
}

func UrlToFeed(url string) (RssFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := httpClient.Get(url)
	if err != nil {
		return RssFeed{}, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return RssFeed{}, err
	}

	rssFeed := RssFeed{}

	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return RssFeed{}, err
	}

	return rssFeed, nil
}
