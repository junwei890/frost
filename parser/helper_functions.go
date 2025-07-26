package parser

import (
	"errors"
	"slices"
	"sort"
	"strings"

	"github.com/junwei890/rumbling/server"
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

type ProcessedText struct {
	Url       string
	Delimited []string
}

func DelimitByPunct(res server.CrawlerRes) (ProcessedText, error) { // delimiting by punctuation to find sentences
	if strings.TrimSpace(res.Doc) == "" {
		return ProcessedText{}, errors.New("unprocessed input")
	}
	delimited := strings.FieldsFunc(res.Doc, func(w rune) bool {
		_, ok := punct[string(w)]
		return ok
	})

	cleanSlice := []string{}
	for _, sent := range delimited {
		cleaned := strings.TrimSpace(sent)
		if cleaned != "" && len(cleaned) > 1 { // filtering out letters and empty space
			cleanSlice = append(cleanSlice, cleaned)
		}
	}

	return ProcessedText{
		Url:       res.URL,
		Delimited: cleanSlice,
	}, nil
}

func DelimitByStop(doc ProcessedText) (ProcessedText, error) { // delimiting by stop words to find phrases
	terms := []string{}
	for _, sent := range doc.Delimited {
		if len(strings.Fields(sent)) > 1 {
			words := strings.Fields(sent)
			curr := 0
			for i, word := range words {
				if _, ok := stopwords[word]; ok {
					phrase := strings.Join(slices.DeleteFunc(words[curr:i], func(w string) bool {
						_, ok := stopwords[w]
						return ok
					}), " ") // dealing with stop words, then joining up the rest to form a phrase
					clean := strings.TrimSpace(phrase)
					if clean != "" {
						terms = append(terms, clean) // case for back to back stop words
					}
					curr = i + 1
				} else if i == len(words)-1 {
					phrase := strings.TrimSpace(strings.Join(words[curr:i+1], " ")) // dealing with the phrase after the last stop word
					terms = append(terms, phrase)
				}
			}
		} else if len(strings.Fields(sent)) == 1 { // appending single words if they aren't stop words
			if _, ok := stopwords[sent]; !ok {
				terms = append(terms, sent)
			}
		} else {
			return ProcessedText{}, errors.New("unprocessed input")
		}
	}
	return ProcessedText{
		Url:       doc.Url,
		Delimited: terms,
	}, nil
}

type CoGraph struct {
	Url   string
	Graph map[string][]string
}

func CoOccurrence(doc ProcessedText) (CoGraph, error) { // a co-occurrence graph is a matrice that shows the frequency of a word's occurrence and co-occurrence with other words
	wordMap := make(map[string][]string)
	for _, term := range doc.Delimited { // creating a map of unique words
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
			return CoGraph{}, errors.New("unprocessed input")
		}
	}

	for _, term := range doc.Delimited { // filling in co-occurrence graph
		if len(strings.Fields(term)) > 1 {
			phrase := strings.Fields(term)
			track := make(map[string]int)
			set := []string{}
			for _, word := range phrase {
				if _, ok := track[word]; !ok {
					track[word] = 1
					set = append(set, word) // creating a set for the phrase
				} else {
					track[word]++ // how many of a word is in a phrase
				}
			}
			for key, value := range track {
				if value == 1 {
					wordMap[key] = slices.Concat(wordMap[key], set) // words cannot co-occur twice with another word in the same phrase
				} else {
					for range value {
						wordMap[key] = append(wordMap[key], key) // words cannot co-occur with themselves if there are more than one of them in the phrase
					}
					for _, word := range set {
						if key != word {
							wordMap[key] = append(wordMap[key], word)
						}
					}
				}
			}
		} else if len(strings.Fields(term)) == 1 {
			wordMap[term] = append(wordMap[term], term) // single words have their presence accounted
		}
	}
	return CoGraph{
		Url:   doc.Url,
		Graph: wordMap,
	}, nil
}

type WordScores struct {
	Url    string
	Scores map[string]float64
}

func DegFreqCalc(graph CoGraph) (WordScores, error) { // word scores are calculated by dividing the degree of a word by its frequency
	scores := make(map[string]float64)
	for key, value := range graph.Graph {
		degree := float64(len(value)) // metric that favors words that occur often as well as within phrases
		if degree == 0.0 {
			return WordScores{}, errors.New("malformed co-occurence table")
		}
		freq := 0.0 // metric that favors words that occur frequently regardless of words which they co-occur with
		for _, word := range value {
			if key == word {
				freq += 1.0
			}
		}
		if freq == 0.0 {
			return WordScores{}, errors.New("malformed co-occurence table")
		}
		scores[key] = degree / freq
	}
	return WordScores{
		Url:    graph.Url,
		Scores: scores,
	}, nil
}

type TermScores struct {
	Url    string
	Scores map[string]float64
}

func TermScoring(scores WordScores, terms ProcessedText) (TermScores, error) { // adding all individual metric scores up per unique term
	if scores.Url != terms.Url {
		return TermScores{}, errors.New("url mismatch")
	}
	track := make(map[string]int) // to check if we've scored words that didn't exist previously
	termScores := make(map[string]float64)
	for _, term := range terms.Delimited {
		if len(strings.Fields(term)) > 1 {
			for word := range strings.FieldsSeq(term) {
				if _, ok := scores.Scores[word]; ok { // checking if word had been scored
					if _, ok := track[word]; !ok {
						track[word] = 1
					}
					termScores[term] += scores.Scores[word]
				} else {
					return TermScores{}, errors.New("word did not exist during scoring")
				}
			}
		} else {
			if _, ok := scores.Scores[term]; ok {
				if _, ok := track[term]; !ok {
					track[term] = 1
				}
				termScores[term] = scores.Scores[term]
			} else {
				return TermScores{}, errors.New("word did not exist during scoring")
			}
		}
	}
	if len(track) != len(scores.Scores) {
		return TermScores{}, errors.New("some scored words were not used")
	}

	return TermScores{
		Url:    scores.Url,
		Scores: termScores,
	}, nil
}

type Keywords struct {
	Url      string
	Keywords []string
}

func Filtering(scores TermScores) Keywords { // according to the paper, we take the top 33% highest scoring terms
	keys := []string{}
	if len(scores.Scores) <= 3 { // returning all keywords if there are only a few
		for key := range scores.Scores {
			keys = append(keys, key)
		}
		return Keywords{
			Url:      scores.Url,
			Keywords: keys,
		}
	}
	for key := range scores.Scores {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool { // sorting the slice of keys corresponding to their value in the map
		return scores.Scores[keys[i]] > scores.Scores[keys[j]]
	})

	taking := len(scores.Scores) / 3
	return Keywords{
		Url:      scores.Url,
		Keywords: keys[0 : taking+1],
	}
}
