package rake

import (
	"reflect"
	"testing"

	"github.com/junwei890/rumbling/server"
)

func TestDelimitByPunct(t *testing.T) {
	testCases := []struct {
		name         string
		input        server.CrawlerRes
		expected     ProcessedText
		errorPresent bool
	}{
		{
			name: "test case 1",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: "good morning, nice weather today!",
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"good morning", "nice weather today"},
			},
			errorPresent: false,
		},
		{
			name: "test case 2",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: "U.S.A, wingstop, basketball?",
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"wingstop", "basketball"},
			},
			errorPresent: false,
		},
		{
			name: "test case 3",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: ",.!?,,..??!!bruh??!!.,",
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"bruh"},
			},
			errorPresent: false,
		},
		{
			name: "test case 4",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: ".,?!",
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{},
			},
			errorPresent: false,
		},
		{
			name: "test case 5",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: "	",
			},
			expected:     ProcessedText{},
			errorPresent: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := DelimitByPunct(testCase.input)
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
		input        ProcessedText
		expected     ProcessedText
		errorPresent bool
	}{
		{
			name: "test case 1",
			input: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"and", "and and", "hello and hi", "and hello", "hi and"},
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"hello", "hi", "hello", "hi"},
			},
			errorPresent: false,
		},
		{
			name: "test case 2",
			input: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"hello hi and hamburgers", "wingstop and fries", "and and fish"},
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"hello hi", "hamburgers", "wingstop", "fries", "fish"},
			},
			errorPresent: false,
		},
		{
			name: "test case 3",
			input: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"and and wow and and", "wow and wow", "and and and and and"},
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"wow", "wow", "wow"},
			},
			errorPresent: false,
		},
		{
			name: "test case 4",
			input: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"hamburgers and crisscut fries", "lemon pepper and lychee"},
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"hamburgers", "crisscut fries", "lemon pepper", "lychee"},
			},
			errorPresent: false,
		},
		{
			name: "test case 5",
			input: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"	"},
			},
			expected:     ProcessedText{},
			errorPresent: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := DelimitByStop(testCase.input)
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
		input        ProcessedText
		expected     CoGraph
		errorPresent bool
	}{
		{
			name: "test case 1",
			input: ProcessedText{
				Url: "bruh",
				Delimited: []string{
					"cloud bill", "sent", "scientific notation",
				},
			},
			expected: CoGraph{
				Url: "bruh",
				Graph: map[string][]string{
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
			input: ProcessedText{
				Url: "",
				Delimited: []string{
					"bubble sort", "quick sort", "ai",
				},
			},
			expected: CoGraph{
				Url: "",
				Graph: map[string][]string{
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
			input: ProcessedText{
				Url: "",
				Delimited: []string{
					"wingstop", "wingstop", "wingstop wingstop",
				},
			},
			expected: CoGraph{
				Url: "",
				Graph: map[string][]string{
					"wingstop": {"wingstop", "wingstop", "wingstop", "wingstop"},
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 4",
			input: ProcessedText{
				Url: "",
				Delimited: []string{
					"term frequency", "inverse document frequency", "tfidf",
				},
			},
			expected: CoGraph{
				Url: "",
				Graph: map[string][]string{
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
			input: ProcessedText{
				Url: "",
				Delimited: []string{
					"stop words", "delimiters", "wing stop", "stop sign",
				},
			},
			expected: CoGraph{
				Url: "",
				Graph: map[string][]string{
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
			input: ProcessedText{
				Url: "",
				Delimited: []string{
					"wingstop", "wingstop", "wingstop bruh wingstop",
				},
			},
			expected: CoGraph{
				Url: "",
				Graph: map[string][]string{
					"wingstop": {"wingstop", "wingstop", "wingstop", "wingstop", "bruh"},
					"bruh":     {"wingstop", "bruh"},
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 7",
			input: ProcessedText{
				Url: "",
				Delimited: []string{
					"wingstop bruh", "bruh wingstop", "wingstop bruh wingstop",
				},
			},
			expected: CoGraph{
				Url: "",
				Graph: map[string][]string{
					"wingstop": {"wingstop", "bruh", "bruh", "wingstop", "wingstop", "wingstop", "bruh"},
					"bruh":     {"wingstop", "bruh", "bruh", "wingstop", "wingstop", "bruh"},
				},
			},
			errorPresent: false,
		},
		{
			name: "test case 8",
			input: ProcessedText{
				Url: "",
				Delimited: []string{
					"hi hi hello hello", "hi", "hello",
				},
			},
			expected: CoGraph{
				Url: "",
				Graph: map[string][]string{
					"hi":    {"hi", "hi", "hello", "hi"},
					"hello": {"hello", "hello", "hi", "hello"},
				},
			},
			errorPresent: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := CoOccurrence(testCase.input)
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
		input        CoGraph
		expected     WordScores
		errorPresent bool
	}{
		{
			name: "test case 1",
			input: CoGraph{
				Url: "bruh",
				Graph: map[string][]string{
					"cloud":      {"cloud", "bill"},
					"bill":       {"cloud", "bill"},
					"sent":       {"sent"},
					"scientific": {"scientific", "notation"},
					"notation":   {"scientific", "notation"},
				},
			},
			expected: WordScores{
				Url: "bruh",
				Scores: map[string]float64{
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
			input: CoGraph{
				Url: "bruh",
				Graph: map[string][]string{
					"stop":       {"stop", "words", "wing", "stop", "stop", "sign"},
					"words":      {"stop", "words"},
					"delimiters": {"delimiters"},
					"wing":       {"wing", "stop"},
					"sign":       {"stop", "sign"},
				},
			},
			expected: WordScores{
				Url: "bruh",
				Scores: map[string]float64{
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
			input: CoGraph{
				Url: "bruh",
				Graph: map[string][]string{
					"term":      {"term", "frequency"},
					"frequency": {"term", "frequency", "inverse", "document", "frequency"},
					"inverse":   {"inverse", "document", "frequency"},
					"document":  {"inverse", "document", "frequency"},
					"tfidf":     {"tfidf"},
				},
			},
			expected: WordScores{
				Url: "bruh",
				Scores: map[string]float64{
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
			input: CoGraph{
				Url: "bruh",
				Graph: map[string][]string{
					"wingstop": {},
				},
			},
			expected:     WordScores{},
			errorPresent: true,
		},
		{
			name: "test case 5",
			input: CoGraph{
				Url: "bruh",
				Graph: map[string][]string{
					"neovim": {"btw"},
				},
			},
			expected:     WordScores{},
			errorPresent: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := DegFreqCalc(testCase.input)
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
		score        WordScores
		terms        ProcessedText
		expected     TermScores
		errorPresent bool
	}{
		{
			name: "test case 1",
			score: WordScores{
				Url:    "bruh",
				Scores: make(map[string]float64),
			},
			terms: ProcessedText{
				Url:       "bruhs",
				Delimited: []string{},
			},
			expected:     TermScores{},
			errorPresent: true,
		},
		{
			name: "test case 2",
			score: WordScores{
				Url: "bruh",
				Scores: map[string]float64{
					"hello": 2.0,
					"hi":    1.5,
					"bye":   3.0,
				},
			},
			terms: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"hello hi bye", "hamburgers"},
			},
			expected:     TermScores{},
			errorPresent: true,
		},
		{
			name: "test case 3",
			score: WordScores{
				Url: "bruh",
				Scores: map[string]float64{
					"hello": 3.0,
					"hi":    1.5,
					"bye":   2.0,
				},
			},
			terms: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"hello and hi", "bye"},
			},
			expected:     TermScores{},
			errorPresent: true,
		},
		{
			name: "test case 4",
			score: WordScores{
				Url: "bruh",
				Scores: map[string]float64{
					"hello": 3.0,
					"hi":    1.5,
					"bye":   2.0,
				},
			},
			terms: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"hello hi"},
			},
			expected:     TermScores{},
			errorPresent: true,
		},
		{
			name: "test case 5",
			score: WordScores{
				Url: "bruh",
				Scores: map[string]float64{
					"buy":      3.4,
					"wingstop": 1.4,
					"today":    2.4,
				},
			},
			terms: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"buy", "wingstop today", "buy wingstop today", "buy today", "buy wingstop"},
			},
			expected: TermScores{
				Url: "bruh",
				Scores: map[string]float64{
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
			result, err := TermScoring(testCase.score, testCase.terms)
			if (err != nil) != testCase.errorPresent {
				t.Errorf("%s failed, expecting err == %v", testCase.name, err)
			} else if comp := reflect.DeepEqual(testCase.expected, result); !comp {
				t.Errorf("%s failed, %v != %v", testCase.name, testCase.expected, result)
			}
		})
	}
}
