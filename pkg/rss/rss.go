package rss

import (
	"encoding/xml"
	storage "goNews/pkg/db"
	"io"
	"net/http"
	"time"
)

type RssNews struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Posts []Post `xml:"item"`
}

// Публикация, получаемая из RSS.
type Post struct {
	Title   string `xml:"title"`       // заголовок публикации
	Content string `xml:"description"` // содержание публикации
	PubTime string `xml:"pubDate"`     // время публикации
	Link    string `xml:"link"`        // ссылка на источник
}

// Получение и парсинг rss потока
func ParseRss(rssUrl string) ([]storage.Post, error) {
	resp, err := http.Get(rssUrl)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	news := new(RssNews)
	err = xml.Unmarshal(body, news)
	if err != nil {
		return nil, err
	}
	var data []storage.Post
	for _, item := range news.Channel.Posts {
		var p storage.Post
		p.Title = item.Title
		p.Content = item.Content
		p.Link = item.Link
		t, err := time.Parse(time.RFC1123, item.PubTime)
		if err != nil {
			return nil, err
		}
		tUnix := t.Unix()
		p.PubTime = tUnix
		data = append(data, p)
	}
	return data, nil

}
