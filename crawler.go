package main

import (
	"context"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/junwei890/rumbling/internal/database"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func (c *crawlerConfig) initCrawl(baseUrl string) {
	c.wg.Add(1)
	go c.crawlPage(baseUrl)
	c.wg.Wait()
}

func (c *crawlerConfig) dataFromHTML(normCurrUrl, htmlBody string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	htmlTree, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return err
	}

	content := []string{}
	for n := range htmlTree.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if urlStruct, err := url.Parse(attr.Val); err != nil {
						return err
					} else if urlStruct.Hostname() == "" {
						c.links[normCurrUrl] = append(c.links[normCurrUrl], c.domain.ResolveReference(urlStruct).String())
					} else {
						c.links[normCurrUrl] = append(c.links[normCurrUrl], attr.Val)
					}
				}
			}
		} else if n.Type == html.ElementNode && n.DataAtom == atom.P {
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				if child.Type == html.TextNode {
					re, err := regexp.Compile(`[^a-zA-Z0-9 .,!?]+`)
					if err != nil {
						return err
					}
					clean := strings.TrimSpace(re.ReplaceAllString(strings.ToLower(child.Data), ""))
					if clean != "" {
						content = append(content, clean)
					}
				}
			}
		}
	}

	clean := strings.TrimSpace(strings.Join(content, " "))
	if clean != "" {
		if err := c.db.InsertData(context.Background(), database.InsertDataParams{
			Url:     normCurrUrl,
			Content: clean,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (c *crawlerConfig) urlVisited(normCurrUrl string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.links[normCurrUrl]; ok {
		return true
	}
	c.links[normCurrUrl] = []string{}
	return false
}

func (c *crawlerConfig) maxReached() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.links) >= c.maxVisits {
		return true
	}
	return false
}

func (c *crawlerConfig) crawlPage(rawCurrUrl string) {
	c.control <- struct{}{}
	defer func() {
		<-c.control
		c.wg.Done()
	}()

	if c.maxReached() {
		return
	}

	currStruct, err := url.Parse(rawCurrUrl)
	if err != nil {
		return
	}
	if c.domain.Hostname() != currStruct.Hostname() {
		return
	}

	normCurrUrl, err := normalizeURL(rawCurrUrl)
	if err != nil {
		return
	}
	if c.urlVisited(normCurrUrl) {
		return
	}

	html, err := getHTML(rawCurrUrl)
	if err != nil {
		return
	}
	if err := c.dataFromHTML(normCurrUrl, html); err != nil {
		return
	}

	for _, link := range c.links[normCurrUrl] {
		c.wg.Add(1)
		log.Printf("crawling %s", link)
		go c.crawlPage(link)
	}
}
