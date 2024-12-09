package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

func (f RSSFeed) String() string {
	stringer := "Channel:\n"
	stringer += fmt.Sprintf("-- Title: %s\n", f.Channel.Title)
	stringer += fmt.Sprintf("-- Link: %s\n", f.Channel.Link)
	stringer += fmt.Sprintf("-- Description: %s\n", f.Channel.Description)

	for _, item := range f.Channel.Item {
		itemStr := item.String()
		stringer += fmt.Sprintf("-- Item:\n")
		for _, line := range strings.Split(itemStr, "\n") {
			stringer += fmt.Sprintf("  %s\n", line)
		}
	}

	return stringer
}

func (f *RSSFeed) Unescape() {
	// removes unwanted escape characters from an RSSFeed
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Description = html.UnescapeString(f.Channel.Description)
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (i RSSItem) String() string {
	stringer := fmt.Sprintf("-- Title: %s\n", i.Title)
	stringer += fmt.Sprintf("-- Link: %s\n", i.Link)
	stringer += fmt.Sprintf("-- Description: %s\n", i.Description)
	stringer += fmt.Sprintf("-- PubDate: %s\n", i.PubDate)
	return stringer
}

func (i *RSSItem) Unescape() {
	// removes unwanted escape characters from an RSSItem
	i.Title = html.UnescapeString(i.Title)
	i.Description = html.UnescapeString(i.Description)
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}

	// create request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %s", err)
	}
	req.Header.Add("User-Agent", "gator")

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %s", err)
	}

	// read request
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %s", err)
	}

	// parse into struct
	var rssFeed RSSFeed
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse into RSSFeed struct: %s", err)
	}

	return &rssFeed, nil
}

func scrapeFeeds(s *state) error {
	background := context.Background()

	// get next feed
	nextFeedToFetch, err := s.db.GetNextFeedToFetch(background)
	if err != nil {
		return fmt.Errorf("Failed to get the next feed to fetch: %s", err)
	}

	// fetch the next feed
	feed, err := fetchFeed(background, nextFeedToFetch.Url)
	if err != nil {
		return fmt.Errorf("Failed to fetch next feed: %s", err)
	}

	// print items in feed
	for _, item := range feed.Channel.Item {
		fmt.Println(item)
	}

	// mark it as fetched
	err = s.db.MarkFeedFetched(background, nextFeedToFetch.ID)
	if err != nil {
		return fmt.Errorf("Failed to mark feed as fetched: %s", err)
	}

	return nil
}
