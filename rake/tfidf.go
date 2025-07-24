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

type Tfidf struct {
	Url         string
	TfidfScores map[string]float64
}

func TfidfCalc(corpus []server.CrawlerRes) []Tfidf {
	agged := []agg{}
	for _, doc := range corpus {
		wordsInDoc := 0
		freqInDoc := make(map[string]int)
		for _, word := range doc.Doc {
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

type rawStats struct {
	url          string
	totalTFIDF   float64
	observations int
	tfidfScores  map[string]float64
}

type Cleaned struct {
	Url   string
	Words []string
}

func StopWordRm(corpus []Tfidf) []Cleaned {
	tfidfStats := []rawStats{}
	for _, doc := range corpus {
		statsPerDoc := rawStats{
			url:          doc.Url,
			totalTFIDF:   0.0,
			observations: 0,
			tfidfScores:  doc.TfidfScores,
		}
		for _, value := range doc.TfidfScores {
			statsPerDoc.totalTFIDF += value
			statsPerDoc.observations++
		}
		tfidfStats = append(tfidfStats, statsPerDoc)
	}

	cleanedDocs := []Cleaned{}
	for _, stats := range tfidfStats {
		cleanedDoc := Cleaned{
			Url:   stats.url,
			Words: []string{},
		}
		summation := 0.0
		for _, value := range stats.tfidfScores {
			summation += math.Pow((value - (stats.totalTFIDF / float64(stats.observations))), 2.0)
		}
		sd := math.Sqrt(summation / float64(stats.observations-1))
		for key, value := range stats.tfidfScores {
			if value < ((stats.totalTFIDF / float64(stats.observations)) - (0.7 * sd)) {
				cleanedDoc.Words = append(cleanedDoc.Words, key)
			}
		}
		cleanedDocs = append(cleanedDocs, cleanedDoc)
	}

	return cleanedDocs
}
