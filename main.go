package main

import (
	"log"
	"os"

	"github.com/junwei890/rumbling/crawler"
	"github.com/junwei890/rumbling/parser"
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
	for _, doc := range res {
		noPunct, err := parser.DelimitByPunct(doc)
		if err != nil {
			log.Fatal(err)
		}
		cleaned, err := parser.DelimitByStop(noPunct)
		if err != err {
			log.Fatal(err)
		}
		coGraph, err := parser.CoOccurrence(cleaned)
		if err != nil {
			log.Fatal(err)
		}
		wordScores, err := parser.DegFreqCalc(coGraph)
		if err != nil {
			log.Fatal(err)
		}
		termScores, err := parser.TermScoring(wordScores, cleaned)
		if err != nil {
			log.Fatal(err)
		}
		filtered := parser.Filtering(termScores)
		for _, keyword := range filtered.Keywords {
			log.Println(keyword)
		}
	}
}
