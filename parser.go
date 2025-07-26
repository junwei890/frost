package main

import "github.com/junwei890/rumbling/internal/database"

func rake(content database.RetrieveDataRow) (keywords, error) {
	byPunct, err := delimitByPunct(content)
	if err != nil {
		return keywords{}, err
	}

	byStop, err := delimitByStop(byPunct)
	if err != nil {
		return keywords{}, err
	}

	graph, err := coOccurrence(byStop)
	if err != nil {
		return keywords{}, err
	}

	wordScore, err := degFreqCalc(graph)
	if err != nil {
		return keywords{}, err
	}

	termScore, err := termScoring(wordScore, byStop)
	if err != nil {
		return keywords{}, err
	}

	return filtering(termScore), nil
}
