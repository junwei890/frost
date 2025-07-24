package rake

import (
	"math"

	"github.com/junwei890/rumbling/server"
)

type agg struct {
	url      string
	docTotal int
	wordFreq map[string]int
}

type tf struct {
	url     string
	tfScore map[string]float64
}

type tfidf struct {
	url        string
	tfidfScore map[string]float64
}

func tfidfCalc(docs []server.CrawlerRes) []tfidf {
	agged := []agg{}
	for _, doc := range docs {
		total := 0
		freq := make(map[string]int)
		for _, word := range doc.Doc {
			total++
			if _, ok := freq[string(word)]; ok {
				freq[string(word)]++
				continue
			}
			freq[string(word)] = 1
		}
		agged = append(agged, agg{
			url:      doc.URL,
			docTotal: total,
			wordFreq: freq,
		})
	}

	tfScores := []tf{}
	idfRef := make(map[string]int)
	for _, data := range agged {
		tfPerDoc := make(map[string]float64)
		for key, value := range data.wordFreq {
			tfPerDoc[key] = float64(value) / float64(data.docTotal)
			if _, ok := idfRef[key]; ok {
				idfRef[key]++
				continue
			}
			idfRef[key] = 1
		}
		tfScores = append(tfScores, tf{
			url:     data.url,
			tfScore: tfPerDoc,
		})
	}

	finalTFIDF := []tfidf{}
	for _, data := range tfScores {
		tfidfPerDoc := make(map[string]float64)
		for key, value := range data.tfScore {
			idfScore := math.Log10(float64(len(docs)) / float64(idfRef[key]))
			tfidfPerDoc[key] = float64(idfScore) * float64(value)
		}
		finalTFIDF = append(finalTFIDF, tfidf{
			url:        data.url,
			tfidfScore: tfidfPerDoc,
		})
	}

	return finalTFIDF
}
