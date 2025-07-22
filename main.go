package main

import (
	"log"
	"os"
	"strconv"

	"github.com/junwei890/frost/scraper"
)

func main() {
	arguments := os.Args
	if len(arguments) < 4 {
		log.Println("no website provided")
		os.Exit(1)
	} else if len(arguments) > 4 {
		log.Println("too many arguments provided")
		os.Exit(1)
	} else {
		log.Printf("starting crawl of %s", arguments[1])
	}

	maxConcurr, err := strconv.Atoi(arguments[2])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	maxSites, err := strconv.Atoi(arguments[3])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	scraper.InitiateCrawl(arguments[1], maxConcurr, maxSites)
}
