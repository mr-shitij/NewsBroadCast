package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	strip "github.com/grokify/html-strip-tags-go"
)

var wg sync.WaitGroup

type News struct {
	Item []item `xml:"channel>item"`
}

type item struct {
	Title       []string `xml:"title"`       // Here we use xml because the output of api / news site is in xml from
	Description []string `xml:"description"` // If the page return output in the form of json then just replace xml with json
	Publish     []string `xml:"pubDate"`
	Link        []string `xml:"link"`
}

type ChannelCarrier struct {
	leabel   string
	Newsdata News
}

func (val item) TruncDESC() string {
	return strip.StripTags(strings.Trim(fmt.Sprint(val.Description), "[]"))
}

func (val item) TruncTILE() string {
	return strings.Trim(fmt.Sprint(val.Title), "[]")
}

func (val item) TruncPUBL() string {
	return strings.Trim(fmt.Sprint(val.Publish), "[]")
}

func (val item) TruncLINK() string {
	return strings.Trim(fmt.Sprint(val.Link), "[]")
}

func GetValByUrl(data chan ChannelCarrier, url string, topic string) {
	defer wg.Done()
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	var Sender ChannelCarrier
	xml.Unmarshal(bytes, &Sender.Newsdata)
	Sender.leabel = topic
	data <- Sender
}

func IndexPage(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello Welcome to frount Page")
}

func NewsPage(writer http.ResponseWriter, request *http.Request) {
	MyUrls := make(map[string]string, 8)
	MyUrls["TopStories"] = "https://timesofindia.indiatimes.com/rssfeedstopstories.cms"
	MyUrls["India"] = "https://timesofindia.indiatimes.com/rssfeeds/-2128936835.cms"
	MyUrls["World"] = "https://timesofindia.indiatimes.com/rssfeeds/296589292.cms"
	MyUrls["Business"] = "https://timesofindia.indiatimes.com/rssfeeds/1898055.cms"
	MyUrls["Sports"] = "https://timesofindia.indiatimes.com/rssfeeds/4719148.cms"
	MyUrls["Science"] = "https://timesofindia.indiatimes.com/rssfeeds/-2128672765.cms"
	MyUrls["Technology"] = "https://gadgets.ndtv.com/rss/feeds"
	MyUrls["Education"] = "https://timesofindia.indiatimes.com/rssfeeds/913168846.cms"

	data := make(chan ChannelCarrier, 1000)
	AllNews := make(map[string]News)
	for topic, url := range MyUrls {
		wg.Add(1)
		go GetValByUrl(data, url, topic)
	}
	wg.Wait()
	close(data)
	for d := range data {
		AllNews[d.leabel] = d.Newsdata
	}
	MyTemplate, _ := template.ParseFiles("Site.html")
	err := MyTemplate.Execute(writer, AllNews)
	if err != nil {
		fmt.Println("\n", err)
	}
}

func main() {
	http.HandleFunc("/", IndexPage)
	http.HandleFunc("/news", NewsPage)
	http.ListenAndServe(":8000", nil)
}
