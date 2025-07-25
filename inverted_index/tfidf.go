package inverted_index

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

type Tfidf struct {
	Url         string
	TfidfScores map[string]float64
}

func TfidfCalc(corpus []server.RakeRes) []Tfidf {
	agged := []agg{}
	for _, doc := range corpus {
		wordsInDoc := 0
		freqInDoc := make(map[string]int)
		for _, word := range doc.Keywords {
			wordsInDoc++
			if _, ok := freqInDoc[string(word)]; ok {
				freqInDoc[string(word)]++
				continue
			}
			freqInDoc[string(word)] = 1
		}
		agged = append(agged, agg{
			url:      doc.URL,
			docTotal: wordsInDoc,
			wordFreq: freqInDoc,
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

	TFIDFScores := []Tfidf{}
	for _, data := range tfScores {
		tfidfPerDoc := make(map[string]float64)
		for key, value := range data.tfScore {
			idfScore := math.Log10(float64(len(corpus)) / float64(idfRef[key]))
			tfidfPerDoc[key] = float64(idfScore) * float64(value)
		}
		TFIDFScores = append(TFIDFScores, Tfidf{
			Url:         data.url,
			TfidfScores: tfidfPerDoc,
		})
	}

	return TFIDFScores
}
