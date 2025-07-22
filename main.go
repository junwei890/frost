package main

import (
	"log"
	"net/url"
	"os"

	"github.com/junwei890/frost/scraper"
)

func main() {
	arguments := os.Args
	if len(arguments) < 2 {
		log.Println("no website provided")
		os.Exit(1)
	} else if len(arguments) > 2 {
		log.Println("too many arguments provided")
		os.Exit(1)
	} else {
		log.Printf("starting crawl of: %s", arguments[1])
	}

	domain, err := url.Parse(arguments[1])
	if err != nil {
		log.Fatal(err)
	}
	pages := make(map[string]int)
	scraper.CrawlPage(domain, arguments[1], pages)

	for key, value := range pages {
		log.Printf("%s: %d", key, value)
	}
}
