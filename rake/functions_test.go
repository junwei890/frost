package rake

import (
	"reflect"
	"testing"

	"github.com/junwei890/rumbling/server"
)

func TestTextProcessing(t *testing.T) {
	testCases := []struct {
		name     string
		input    server.CrawlerRes
		expected ProcessedText
	}{
		{
			name: "test case 1",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: []string{"i", "love", "pizza", "and", "hamburgers"},
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"love pizza", "hamburgers"},
			},
		},
		{
			name: "test case 2",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: []string{"i", "need", "wingstop", "again"},
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"need wingstop"},
			},
		},
		{
			name: "test case 3",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: []string{"i", "i", "i", "i", "i"},
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{},
			},
		},
		{
			name: "test case 4",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: []string{"hello", "nice", "meeting", "you", "again"},
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"hello nice meeting"},
			},
		},
		{
			name: "test case 5",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: []string{"i", "i", "hello", "i", "i", "hi", "i"},
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"hello", "hi"},
			},
		},
		{
			name: "test case 6",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: []string{"hello", "good", "morning", "wonderful", "day"},
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{"hello good morning wonderful day"},
			},
		},
		{
			name: "test case 7",
			input: server.CrawlerRes{
				URL: "bruh",
				Doc: []string{},
			},
			expected: ProcessedText{
				Url:       "bruh",
				Delimited: []string{},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			output := TextProcessing(testCase.input)
			if comp := reflect.DeepEqual(output, testCase.expected); !comp {
				t.Errorf("%s failed, %v != %v", testCase.name, output.Delimited, testCase.expected.Delimited)
			}
		})
	}
}

func TestCoOccurence(t *testing.T) {
	testCases := []struct {
		name     string
		input    ProcessedText
		expected CoGraph
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
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := CoOccurence(testCase.input)
			if comp := reflect.DeepEqual(result, testCase.expected); !comp {
				t.Errorf("%s failed, %v != %v", testCase.name, result.Graph, testCase.expected.Graph)
			}
		})
	}
}

func TestDegFreqCalc(t *testing.T) {
	testCases := []struct {
		name     string
		input    CoGraph
		expected WordScores
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
					"cloud": 2.0,
					"bill": 2.0,
					"sent": 1.0,
					"scientific": 2.0,
					"notation": 2.0,
				},
			},
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
					"stop": 2.0,
					"words": 2.0,
					"delimiters": 1.0,
					"wing": 2.0,
					"sign": 2.0,
				},
			},
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
					"term": 2.0,
					"frequency": 2.5,
					"inverse": 3.0,
					"document": 3.0,
					"tfidf": 1.0,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := DegFreqCalc(testCase.input)
			if comp := reflect.DeepEqual(result.Scores, testCase.expected.Scores); !comp {
				t.Errorf("%s failed, %v != %v", testCase.name, result.Scores, testCase.expected)
			}
		})
	}
}
