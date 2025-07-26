package main

import (
	"reflect"
	"testing"

	"github.com/junwei890/rumbling/internal/database"
)

func TestDelimitByPunct(t *testing.T) {
	testCases := []struct {
		name         string
		input        database.RetrieveDataRow
		expected     processedText
		errorPresent bool
	}{
		{
			name: "test case 1",
			input: database.RetrieveDataRow{
				Url:     "bruh",
				Content: "good morning, nice weather today!",
			},
			expected: processedText{
				url:       "bruh",
				delimited: []string{"good morning", "nice weather today"},
			},
			errorPresent: false,
		},
		{
			name: "test case 2",
			input: database.RetrieveDataRow{
				Url:     "bruh",
				Content: "U.S.A, wingstop, basketball?",
			},
			expected: processedText{
				url:       "bruh",
				delimited: []string{"wingstop", "basketball"},
			},
			errorPresent: false,
		},
		{
			name: "test case 3",
			input: database.RetrieveDataRow{
				Url:     "bruh",
				Content: ",.!?,,..??!!bruh??!!.,",
			},
			expected: processedText{
				url:       "bruh",
				delimited: []string{"bruh"},
			},
			errorPresent: false,
		},
		{
			name: "test case 4",
			input: database.RetrieveDataRow{
				Url:     "bruh",
				Content: ".,?!",
			},
			expected: processedText{
				url:       "bruh",
				delimited: []string{},
			},
			errorPresent: false,
		},
		{
			name: "test case 5",
			input: database.RetrieveDataRow{
				Url:     "bruh",
				Content: "	",
			},
			expected:     processedText{},
			errorPresent: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := delimitByPunct(testCase.input)
			if (err != nil) != testCase.errorPresent {
				t.Errorf("%s failed, expecting err = %v", testCase.name, err)
			} else if comp := reflect.DeepEqual(result, testCase.expected); !comp {
				t.Errorf("%s failed, %v != %v", testCase.name, result, testCase.expected)
			}
		})
	}
}

func TestDelimitByStop(t *testing.T) {
	testCases := []struct {
		name         string
		input        processedText
		expected     processedText
		errorPresent bool
	}{
		{
			name: "test case 1",
			input: processedText{
				url:       "bruh",
				delimited: []string{"and", "and and", "hello and hi", "and hello", "hi and"},
			},
			expected: processedText{
				url:       "bruh",
				delimited: []string{"hello", "hi", "hello", "hi"},
			},
			errorPresent: false,
		},
		{
			name: "test case 2",
			input: processedText{
				url:       "bruh",
				delimited: []string{"hello hi and hamburgers", "wingstop and fries", "and and fish"},
			},
			expected: processedText{
				url:       "bruh",
				delimited: []string{"hello hi", "hamburgers", "wingstop", "fries", "fish"},
			},
			errorPresent: false,
		},
		{
			name: "test case 3",
			input: processedText{
				url:       "bruh",
				delimited: []string{"and and wow and and", "wow and wow", "and and and and and"},
			},
			expected: processedText{
				url:       "bruh",
				delimited: []string{"wow", "wow", "wow"},
			},
			errorPresent: false,
		},
		{
			name: "test case 4",
			input: processedText{
				url:       "bruh",
				delimited: []string{"hamburgers and crisscut fries", "lemon pepper and lychee"},
			},
			expected: processedText{
				url:       "bruh",
				delimited: []string{"hamburgers", "crisscut fries", "lemon pepper", "lychee"},
			},
			errorPresent: false,
		},
		{
			name: "test case 5",
			input: processedText{
				url:       "bruh",
				delimited: []string{"	"},
			},
			expected:     processedText{},
			errorPresent: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := delimitByStop(testCase.input)
			if (err != nil) != testCase.errorPresent {
				t.Errorf("%s failed, expected err = %v", testCase.name, err)
			} else if comp := reflect.DeepEqual(result, testCase.expected); !comp {
				t.Errorf("%s failed, %v != %v", testCase.name, result, testCase.expected)
			}
		})
	}
}

func TestCoOccurrence(t *testing.T) {
	testCases := []struct {
		name         string
		input        processedText
		expected     coGraph
		errorPresent bool
	}{
		{
			name: "test case 1",
			input: processedText{
				url: "bruh",
				delimited: []string{
					"cloud bill", "sent", "scientific notation",
				},
			},
			expected: coGraph{
				url: "bruh",
				graph: map[string][]string{
					"cloud":      {"cloud", "bill"},
					"bill":       {"cloud", "bill"},
					"sent":       {"sent"},
					"scientific": {"scientific", "notation"},
					"notation":   {"scientific", "notation"},
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 2",
			input: processedText{
				url: "",
				delimited: []string{
					"bubble sort", "quick sort", "ai",
				},
			},
			expected: coGraph{
				url: "",
				graph: map[string][]string{
					"bubble": {"bubble", "sort"},
					"sort":   {"bubble", "sort", "quick", "sort"},
					"quick":  {"quick", "sort"},
					"ai":     {"ai"},
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 3",
			input: processedText{
				url: "",
				delimited: []string{
					"wingstop", "wingstop", "wingstop wingstop",
				},
			},
			expected: coGraph{
				url: "",
				graph: map[string][]string{
					"wingstop": {"wingstop", "wingstop", "wingstop", "wingstop"},
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 4",
			input: processedText{
				url: "",
				delimited: []string{
					"term frequency", "inverse document frequency", "tfidf",
				},
			},
			expected: coGraph{
				url: "",
				graph: map[string][]string{
					"term":      {"term", "frequency"},
					"frequency": {"term", "frequency", "inverse", "document", "frequency"},
					"inverse":   {"inverse", "document", "frequency"},
					"document":  {"inverse", "document", "frequency"},
					"tfidf":     {"tfidf"},
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 5",
			input: processedText{
				url: "",
				delimited: []string{
					"stop words", "delimiters", "wing stop", "stop sign",
				},
			},
			expected: coGraph{
				url: "",
				graph: map[string][]string{
					"stop":       {"stop", "words", "wing", "stop", "stop", "sign"},
					"words":      {"stop", "words"},
					"delimiters": {"delimiters"},
					"wing":       {"wing", "stop"},
					"sign":       {"stop", "sign"},
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 6",
			input: processedText{
				url: "",
				delimited: []string{
					"wingstop", "wingstop", "wingstop bruh wingstop",
				},
			},
			expected: coGraph{
				url: "",
				graph: map[string][]string{
					"wingstop": {"wingstop", "wingstop", "wingstop", "wingstop", "bruh"},
					"bruh":     {"wingstop", "bruh"},
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 7",
			input: processedText{
				url: "",
				delimited: []string{
					"wingstop bruh", "bruh wingstop", "wingstop bruh wingstop",
				},
			},
			expected: coGraph{
				url: "",
				graph: map[string][]string{
					"wingstop": {"wingstop", "bruh", "bruh", "wingstop", "wingstop", "wingstop", "bruh"},
					"bruh":     {"wingstop", "bruh", "bruh", "wingstop", "wingstop", "bruh"},
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 8",
			input: processedText{
				url: "",
				delimited: []string{
					"hi hi hello hello", "hi", "hello",
				},
			},
			expected: coGraph{
				url: "",
				graph: map[string][]string{
					"hi":    {"hi", "hi", "hello", "hi"},
					"hello": {"hello", "hello", "hi", "hello"},
				},
			},
			errorPresent: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := coOccurrence(testCase.input)
			if (err != nil) != testCase.errorPresent {
				t.Errorf("%s failed, expecting err = %v", testCase.name, err)
			} else if comp := reflect.DeepEqual(result, testCase.expected); !comp {
				t.Errorf("%s failed, %v != %v", testCase.name, result, testCase.expected)
			}
		})
	}
}

func TestDegFreqCalc(t *testing.T) {
	testCases := []struct {
		name         string
		input        coGraph
		expected     wordScores
		errorPresent bool
	}{
		{
			name: "test case 1",
			input: coGraph{
				url: "bruh",
				graph: map[string][]string{
					"cloud":      {"cloud", "bill"},
					"bill":       {"cloud", "bill"},
					"sent":       {"sent"},
					"scientific": {"scientific", "notation"},
					"notation":   {"scientific", "notation"},
				},
			},
			expected: wordScores{
				url: "bruh",
				scores: map[string]float64{
					"cloud":      2.0,
					"bill":       2.0,
					"sent":       1.0,
					"scientific": 2.0,
					"notation":   2.0,
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 2",
			input: coGraph{
				url: "bruh",
				graph: map[string][]string{
					"stop":       {"stop", "words", "wing", "stop", "stop", "sign"},
					"words":      {"stop", "words"},
					"delimiters": {"delimiters"},
					"wing":       {"wing", "stop"},
					"sign":       {"stop", "sign"},
				},
			},
			expected: wordScores{
				url: "bruh",
				scores: map[string]float64{
					"stop":       2.0,
					"words":      2.0,
					"delimiters": 1.0,
					"wing":       2.0,
					"sign":       2.0,
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 3",
			input: coGraph{
				url: "bruh",
				graph: map[string][]string{
					"term":      {"term", "frequency"},
					"frequency": {"term", "frequency", "inverse", "document", "frequency"},
					"inverse":   {"inverse", "document", "frequency"},
					"document":  {"inverse", "document", "frequency"},
					"tfidf":     {"tfidf"},
				},
			},
			expected: wordScores{
				url: "bruh",
				scores: map[string]float64{
					"term":      2.0,
					"frequency": 2.5,
					"inverse":   3.0,
					"document":  3.0,
					"tfidf":     1.0,
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 4",
			input: coGraph{
				url: "bruh",
				graph: map[string][]string{
					"wingstop": {},
				},
			},
			expected:     wordScores{},
			errorPresent: true,
		},
		{
			name: "test case 5",
			input: coGraph{
				url: "bruh",
				graph: map[string][]string{
					"neovim": {"btw"},
				},
			},
			expected:     wordScores{},
			errorPresent: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := degFreqCalc(testCase.input)
			if (err != nil) != testCase.errorPresent {
				t.Errorf("%s failed, expecting err == %v", testCase.name, err)
			} else if comp := reflect.DeepEqual(result, testCase.expected); !comp {
				t.Errorf("%s failed, %v != %v", testCase.name, result, testCase.expected)
			}
		})
	}
}

func TestTermScoring(t *testing.T) {
	testCases := []struct {
		name         string
		score        wordScores
		terms        processedText
		expected     termScores
		errorPresent bool
	}{
		{
			name: "test case 1",
			score: wordScores{
				url:    "bruh",
				scores: make(map[string]float64),
			},
			terms: processedText{
				url:       "bruhs",
				delimited: []string{},
			},
			expected:     termScores{},
			errorPresent: true,
		},
		{
			name: "test case 2",
			score: wordScores{
				url: "bruh",
				scores: map[string]float64{
					"hello": 2.0,
					"hi":    1.5,
					"bye":   3.0,
				},
			},
			terms: processedText{
				url:       "bruh",
				delimited: []string{"hello hi bye", "hamburgers"},
			},
			expected:     termScores{},
			errorPresent: true,
		},
		{
			name: "test case 3",
			score: wordScores{
				url: "bruh",
				scores: map[string]float64{
					"hello": 3.0,
					"hi":    1.5,
					"bye":   2.0,
				},
			},
			terms: processedText{
				url:       "bruh",
				delimited: []string{"hello and hi", "bye"},
			},
			expected:     termScores{},
			errorPresent: true,
		},
		{
			name: "test case 4",
			score: wordScores{
				url: "bruh",
				scores: map[string]float64{
					"hello": 3.0,
					"hi":    1.5,
					"bye":   2.0,
				},
			},
			terms: processedText{
				url:       "bruh",
				delimited: []string{"hello hi"},
			},
			expected:     termScores{},
			errorPresent: true,
		},
		{
			name: "test case 5",
			score: wordScores{
				url: "bruh",
				scores: map[string]float64{
					"buy":      3.4,
					"wingstop": 1.4,
					"today":    2.4,
				},
			},
			terms: processedText{
				url:       "bruh",
				delimited: []string{"buy", "wingstop today", "buy wingstop today", "buy today", "buy wingstop"},
			},
			expected: termScores{
				url: "bruh",
				scores: map[string]float64{
					"buy":                3.4,
					"wingstop today":     3.8,
					"buy wingstop today": 7.199999999999999,
					"buy today":          5.8,
					"buy wingstop":       4.8,
				},
			},
			errorPresent: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := termScoring(testCase.score, testCase.terms)
			if (err != nil) != testCase.errorPresent {
				t.Errorf("%s failed, expecting err == %v", testCase.name, err)
			} else if comp := reflect.DeepEqual(testCase.expected, result); !comp {
				t.Errorf("%s failed, %v != %v", testCase.name, testCase.expected, result)
			}
		})
	}
}

func TestFiltering(t *testing.T) {
	testCases := []struct {
		name     string
		input    termScores
		expected keywords
	}{
		{
			name: "test case 1",
			input: termScores{
				url: "bruh",
				scores: map[string]float64{
					"wingstop good":  5.0,
					"monday tuesday": 1.5,
					"wow":            6.0,
					"terrific":       4.5,
					"awesome job":    3.0,
					"hello":          4.0,
				},
			},
			expected: keywords{
				url:      "bruh",
				keywords: []string{"wow", "wingstop good", "terrific"},
			},
		},
		{
			name: "test case 2",
			input: termScores{
				url: "bruh",
				scores: map[string]float64{
					"hello world": 4.5,
					"npm install": 3.4,
					"hello":       3.2,
					"good day":    1.5,
				},
			},
			expected: keywords{
				url:      "bruh",
				keywords: []string{"hello world", "npm install"},
			},
		},
		{
			name: "test case 3",
			input: termScores{
				url: "bruh",
				scores: map[string]float64{
					"hello world": 4.5,
					"golang":      3.0,
					"wingstop":    2.0,
				},
			},
			expected: keywords{
				url:      "bruh",
				keywords: []string{"hello world", "golang", "wingstop"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := filtering(testCase.input)
			if comp := reflect.DeepEqual(result, testCase.expected); !comp {
				t.Errorf("%s failed, %v != %v", testCase.name, result, testCase.expected)
			}
		})
	}
}
