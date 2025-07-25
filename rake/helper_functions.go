package rake

import (
	"slices"
	"strings"

	"github.com/junwei890/rumbling/server"
)

var delimiters = map[string]struct{}{
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

type ProcessedText struct {
	Url       string
	Delimited []string
}

func TextProcessing(doc server.CrawlerRes) ProcessedText {
	curr := 0
	cleanedDoc := []string{}
	for i, word := range doc.Doc {
		if _, ok := delimiters[word]; ok {
			termSlice := slices.DeleteFunc(doc.Doc[curr:i], func(word string) bool {
				_, ok := delimiters[word]
				return ok
			})
			term := strings.Join(termSlice, " ")
			if strings.TrimSpace(term) != "" {
				cleanedDoc = append(cleanedDoc, strings.TrimSpace(term))
				curr = i + 1
			}
		} else if i == (len(doc.Doc) - 1) {
			term := strings.Join(doc.Doc[curr:i+1], " ")
			if strings.TrimSpace(term) == "" {
				continue
			}
			cleanedDoc = append(cleanedDoc, strings.TrimSpace(term))
		}
	}
	return ProcessedText{
		Url:       doc.URL,
		Delimited: cleanedDoc,
	}
}

type CoGraph struct {
	Url   string
	Graph map[string][]string
}

func CoOccurence(doc ProcessedText) CoGraph {
	wordMap := make(map[string][]string)
	for _, term := range doc.Delimited {
		if len(strings.Fields(term)) > 1 {
			for word := range strings.FieldsSeq(term) {
				if _, ok := wordMap[word]; !ok {
					wordMap[word] = []string{}
				}
			}
		} else if _, ok := wordMap[term]; !ok {
			wordMap[term] = []string{}
		}
	}

	for _, term := range doc.Delimited {
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
					wordMap[key] = slices.Concat(wordMap[key], set)
				} else {
					for range value {
						wordMap[key] = append(wordMap[key], key)
					}
					for _, word := range set {
						if key != word {
							wordMap[key] = append(wordMap[key], word)
						}
					}
				}
			}
		} else {
			wordMap[term] = append(wordMap[term], term)
		}
	}
	return CoGraph{
		Url:   doc.Url,
		Graph: wordMap,
	}
}
