package main

import (
	"log"
	"os"

	"github.com/junwei890/rumbling/crawler"
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
		log.Printf("starting crawl of %s", arguments[1])
	}

	res, err := crawler.InitiateCrawl(arguments[1])
	if err != nil {
		log.Fatal(err)
	}

	for _, crawl := range res {
		log.Println(crawl.URL)
		log.Println(crawl.Doc)
	}
}
