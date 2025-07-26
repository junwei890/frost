package main

import (
	"errors"
	"slices"
	"sort"
	"strings"

	"github.com/junwei890/rumbling/internal/database"
)

var stopwords = map[string]struct{}{
	"i": {}, "im": {}, "ive": {}, "ill": {}, "id": {}, "me": {}, "my": {}, "myself": {}, "we": {}, "wed": {},
	"were": {}, "weve": {}, "our": {}, "ours": {}, "ourselves": {}, "you": {}, "youre": {}, "youve": {},
	"youll": {}, "youd": {}, "your": {}, "yours": {}, "yourself": {}, "yourselves": {}, "he": {}, "hed": {},
	"hell": {}, "hes": {}, "him": {}, "his": {}, "himself": {}, "she": {}, "shed": {}, "shell": {}, "shes": {},
	"her": {}, "hers": {}, "herself": {}, "it": {}, "itd": {}, "itll": {}, "its": {}, "itself": {},
	"they": {}, "theyd": {}, "theyll": {}, "theyre": {}, "theyve": {}, "them": {}, "their": {}, "theirs": {},
	"themselves": {}, "what": {}, "whats": {}, "which": {}, "who": {}, "whos": {}, "whom": {}, "this": {},
	"that": {}, "thats": {}, "these": {}, "those": {}, "am": {}, "is": {}, "are": {}, "was": {},
	"be": {}, "been": {}, "being": {}, "have": {}, "has": {}, "had": {}, "having": {}, "do": {}, "does": {},
	"did": {}, "doing": {}, "a": {}, "an": {}, "the": {}, "and": {}, "but": {}, "if": {}, "or": {}, "because": {},
	"as": {}, "until": {}, "while": {}, "of": {}, "at": {}, "by": {}, "for": {}, "with": {}, "about": {},
	"against": {}, "between": {}, "into": {}, "through": {}, "during": {}, "before": {}, "after": {}, "above": {},
	"below": {}, "to": {}, "from": {}, "up": {}, "down": {}, "in": {}, "out": {}, "on": {}, "off": {}, "over": {},
	"under": {}, "again": {}, "further": {}, "then": {}, "once": {}, "here": {}, "there": {}, "when": {}, "where": {},
	"why": {}, "how": {}, "all": {}, "any": {}, "both": {}, "each": {}, "few": {}, "more": {}, "most": {},
	"other": {}, "some": {}, "such": {}, "no": {}, "nor": {}, "not": {}, "only": {}, "own": {}, "same": {},
	"so": {}, "than": {}, "too": {}, "very": {}, "can": {}, "will": {}, "just": {}, "dont": {}, "doesnt": {},
	"didnt": {}, "hasnt": {}, "havent": {}, "isnt": {}, "wasnt": {}, "wont": {}, "would": {}, "wouldnt": {},
	"could": {}, "couldnt": {}, "should": {}, "shouldnt": {}, "must": {}, "mustnt": {}, "let": {}, "lets": {},
	"theres": {}, "wouldve": {}, "couldve": {}, "shouldve": {}, "s": {}, "t": {}, "don": {}, "now": {},
}

var punct = map[string]struct{}{
	".": {}, ",": {}, "?": {}, "!": {},
}

type processedText struct {
	url       string
	delimited []string
}

func delimitByPunct(res database.RetrieveDataRow) (processedText, error) { // delimiting by punctuation to find sentences
	if strings.TrimSpace(res.Content) == "" {
		return processedText{}, errors.New("unprocessed input")
	}
	delimited := strings.FieldsFunc(res.Content, func(w rune) bool {
		_, ok := punct[string(w)]
		return ok
	})

	cleanSlice := []string{}
	for _, sent := range delimited {
		cleaned := strings.TrimSpace(sent)
		if cleaned != "" && len(cleaned) > 1 {
			cleanSlice = append(cleanSlice, cleaned)
		}
	}

	return processedText{
		url:       res.Url,
		delimited: cleanSlice,
	}, nil
}

func delimitByStop(doc processedText) (processedText, error) { // delimiting by stop words to find phrases
	terms := []string{}
	for _, sent := range doc.delimited {
		if len(strings.Fields(sent)) > 1 {
			words := strings.Fields(sent)
			curr := 0
			for i, word := range words {
				if _, ok := stopwords[word]; ok {
					phrase := strings.Join(slices.DeleteFunc(words[curr:i], func(w string) bool {
						_, ok := stopwords[w]
						return ok
					}), " ") // joining up words between 2 stop words
					clean := strings.TrimSpace(phrase)
					if clean != "" {
						terms = append(terms, clean) // case for back to back stop words
					}
					curr = i + 1
				} else if i == len(words)-1 {
					phrase := strings.TrimSpace(strings.Join(words[curr:i+1], " ")) // dealing with a potential phrase after the last stop word
					terms = append(terms, phrase)
				}
			}
		} else if len(strings.Fields(sent)) == 1 {
			if _, ok := stopwords[sent]; !ok {
				terms = append(terms, sent)
			}
		} else {
			return processedText{}, errors.New("unprocessed input")
		}
	}
	return processedText{
		url:       doc.url,
		delimited: terms,
	}, nil
}

type coGraph struct {
	url   string
	graph map[string][]string
}

func coOccurrence(doc processedText) (coGraph, error) { // plotting the co-occurence graph to easily calculate word scores
	wordMap := make(map[string][]string)
	for _, term := range doc.delimited {
		if len(strings.Fields(term)) > 1 {
			for word := range strings.FieldsSeq(term) {
				if _, ok := wordMap[word]; !ok {
					wordMap[word] = []string{}
				}
			}
		} else if len(strings.Fields(term)) == 1 {
			if _, ok := wordMap[term]; !ok {
				wordMap[term] = []string{}
			}
		} else {
			return coGraph{}, errors.New("unprocessed input")
		}
	}

	for _, term := range doc.delimited {
		if len(strings.Fields(term)) > 1 {
			phrase := strings.Fields(term)
			track := make(map[string]int)
			set := []string{}
			for _, word := range phrase {
				if _, ok := track[word]; !ok {
					track[word] = 1
					set = append(set, word)
				} else {
					track[word]++
				}
			}
			for key, value := range track {
				if value == 1 {
					wordMap[key] = slices.Concat(wordMap[key], set) // words cannot co-occur twice with another word in the same phrase if there are 2 of said word
				} else {
					for range value {
						wordMap[key] = append(wordMap[key], key) // words cannot co-occur with themselves if there are more than one of them in the same phrase
					}
					for _, word := range set {
						if key != word {
							wordMap[key] = append(wordMap[key], word)
						}
					}
				}
			}
		} else if len(strings.Fields(term)) == 1 {
			wordMap[term] = append(wordMap[term], term)
		}
	}
	return coGraph{
		url:   doc.url,
		graph: wordMap,
	}, nil
}

type wordScores struct {
	url    string
	scores map[string]float64
}

func degFreqCalc(graph coGraph) (wordScores, error) { // word scores = deg/freq
	scores := make(map[string]float64)
	for key, value := range graph.graph {
		degree := float64(len(value))
		if degree == 0.0 {
			return wordScores{}, errors.New("malformed co-occurence table")
		}
		freq := 0.0
		for _, word := range value {
			if key == word {
				freq += 1.0
			}
		}
		if freq == 0.0 {
			return wordScores{}, errors.New("malformed co-occurence table")
		}
		scores[key] = degree / freq
	}
	return wordScores{
		url:    graph.url,
		scores: scores,
	}, nil
}

type termScores struct {
	url    string
	scores map[string]float64
}

func termScoring(scores wordScores, terms processedText) (termScores, error) { // adding up all the scores for words that make up a phrase
	if scores.url != terms.url {
		return termScores{}, errors.New("url mismatch")
	}
	track := make(map[string]int)
	termScore := make(map[string]float64)
	for _, term := range terms.delimited {
		if len(strings.Fields(term)) > 1 {
			for word := range strings.FieldsSeq(term) {
				if _, ok := scores.scores[word]; ok {
					if _, ok := track[word]; !ok {
						track[word] = 1
					}
					termScore[term] += scores.scores[word]
				} else {
					return termScores{}, errors.New("word did not exist during scoring")
				}
			}
		} else {
			if _, ok := scores.scores[term]; ok {
				if _, ok := track[term]; !ok {
					track[term] = 1
				}
				termScore[term] = scores.scores[term]
			} else {
				return termScores{}, errors.New("word did not exist during scoring")
			}
		}
	}
	if len(track) != len(scores.scores) {
		return termScores{}, errors.New("some scored words were not used")
	}

	return termScores{
		url:    scores.url,
		scores: termScore,
	}, nil
}

type keywords struct {
	url      string
	keywords []string
}

func filtering(scores termScores) keywords { // we take the top 33% highest scoring terms
	keys := []string{}
	if len(scores.scores) <= 3 {
		for key := range scores.scores {
			keys = append(keys, key)
		}
		return keywords{
			url:      scores.url,
			keywords: keys,
		}
	}
	for key := range scores.scores {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return scores.scores[keys[i]] > scores.scores[keys[j]]
	})

	taking := len(scores.scores) / 3
	return keywords{
		url:      scores.url,
		keywords: keys[0 : taking+1],
	}
}
