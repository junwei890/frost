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
				Url:     "bruh",
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
				Url:     "bruh",
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
				Url:     "bruh",
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
				Url:     "bruh",
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
				Url:     "bruh",
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
				Url:     "bruh",
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
