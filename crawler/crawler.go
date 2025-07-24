package crawler

import (
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/junwei890/rumbling/server"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type data struct {
	title string
	text  []string
	links []string
}

type config struct {
	metadata  map[string]data
	domain    *url.URL
	mu        *sync.Mutex
	wg        *sync.WaitGroup
	control   chan struct{}
	maxVisits int
}

func InitiateCrawl(baseURL string) ([]server.CrawlerRes, error) {
	domain, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	local := config{
		metadata:  make(map[string]data),
		domain:    domain,
		mu:        &sync.Mutex{},
		wg:        &sync.WaitGroup{},
		control:   make(chan struct{}, 5),
		maxVisits: 20,
	}

	local.wg.Add(1)
	go local.crawlPage(baseURL)
	local.wg.Wait()

	res := []server.CrawlerRes{}
	for key, value := range local.metadata {
		if value.text == nil || value.title == "" {
			continue
		}
		res = append(res, server.CrawlerRes{
			URL:     key,
			Title:   value.title,
			Content: value.text,
		})
	}
	return res, nil
}

func (c *config) dataFromHTML(normCurrURL, htmlBody string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	htmlTree, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return err
	}
	urlData := data{
		title: "",
		text:  []string{},
		links: []string{},
	}

	for n := range htmlTree.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if urlStruct, err := url.Parse(attr.Val); err != nil {
						return err
					} else if urlStruct.Hostname() == "" {
						urlData.links = append(urlData.links, c.domain.ResolveReference(urlStruct).String())
						continue
					}
					urlData.links = append(urlData.links, attr.Val)
				}
			}
		}
		if n.Type == html.ElementNode && n.DataAtom == atom.Title {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					urlData.title = strings.ToLower(strings.Join(strings.Fields(c.Data), " "))
				}
			}
		}
		if n.Type == html.ElementNode && (n.DataAtom == atom.P || n.DataAtom == atom.H1) {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					urlData.text = append(urlData.text, strings.ToLower(strings.Join(strings.Fields(c.Data), " ")))
				}
			}
		}
	}
	c.metadata[normCurrURL] = urlData
	return nil
}

func (c *config) urlVisited(normCurrURL string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.metadata[normCurrURL]; ok {
		return true
	}
	c.metadata[normCurrURL] = data{}
	return false
}

func (c *config) maxReached() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.metadata) >= c.maxVisits {
		return true
	}
	return false
}

func (c *config) crawlPage(rawCurrURL string) {
	c.control <- struct{}{}
	defer func() {
		<-c.control
		c.wg.Done()
	}()

	if c.maxReached() {
		return
	}
	currStruct, err := url.Parse(rawCurrURL)
	if err != nil {
		return
	}
	if c.domain.Hostname() != currStruct.Hostname() {
		return
	}
	normCurrURL, err := normalizeURL(rawCurrURL)
	if err != nil {
		return
	}
	if c.urlVisited(normCurrURL) {
		return
	}

	html, err := getHTML(rawCurrURL)
	if err != nil {
		return
	}
	if err := c.dataFromHTML(normCurrURL, html); err != nil {
		return
	}

	for _, link := range c.metadata[normCurrURL].links {
		c.wg.Add(1)
		log.Printf("crawling %s", link)
		go c.crawlPage(link)
	}
}
