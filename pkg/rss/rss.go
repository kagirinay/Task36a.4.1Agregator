package rss

import (
	"Task36a.4.1Aggregator/pkg/storage"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Item Структура для отдельного поста.
type Item struct {
	Title       string        `xml:"title"`
	Content     template.HTML `xml:"content"`
	PublishedAt string        `xml:"publishedAt"`
	Link        string        `xml:"link"`
}

// MyXMLStruct Структура данных, получаемая из RSS.
type MyXMLStruct struct {
	ItemList []Item `xml:"channel>item"`
}

// RSSToStruct Преобразование полученных XML данных в заданную структуру, затем в массив новостей.
func RSSToStruct(link string) ([]storage.Post, error) {
	var posts MyXMLStruct
	if xmlBytes, err := receivingXML(link) {
		if err != nil {
			log.Printf("Ошибка при получении данных из XML: %v", err)
		} else {
			xml.Unmarshal(xmlBytes, &posts)
		}
	}
	var news []storage.Post
	for j := range posts.ItemList {
		var item storage.Post
		item.Title = posts.ItemList[j].Title
		item.Content = string(posts.ItemList[j].Content)
		item.Link = posts.ItemList[j].Link

		posts.ItemList[j].PublishedAt = strings.ReplaceAll(posts.ItemList[j].PublishedAt, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 0700", posts.ItemList[j].PublishedAt)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", posts.ItemList[j].PublishedAt)
		}
		if err == nil {
			item.PublishedAt = t.Unix()
		}
		news = append(news, item)
	}
	return news, nil
}

// Получение XML данных по ссылке.
func receivingXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}
	return data, nil
}
