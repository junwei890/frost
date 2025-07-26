package main

import (
	"log"
	"os"

	"github.com/junwei890/rumbling/crawler"
	"github.com/junwei890/rumbling/rake"
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

	for _, res := range res {
		doc := rake.TextProcessing(res)
		graph := rake.CoOccurence(doc)
		scores, err := rake.DegFreqCalc(graph)
		if err != nil {
			log.Fatal(err)
		}
		termscores, err := rake.TermScoring(scores, doc)
		for key, value := range termscores.Scores {
			log.Printf("%s: %.2f", key, value)
		}
	}
}
