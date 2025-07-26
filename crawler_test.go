package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	testCases := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "test case 1",
			url:      "http://www.hello.com/world",
			expected: "www.hello.com/world",
		},
		{
			name:     "test case 2",
			url:      "http://www.hello.com/world/",
			expected: "www.hello.com/world",
		},
		{
			name:     "test case 3",
			url:      "https://www.hello.com/world",
			expected: "www.hello.com/world",
		},
		{
			name:     "test case 4",
			url:      "https://www.hello.com/world/",
			expected: "www.hello.com/world",
		},
		{
			name:     "test case 5",
			url:      "https://www.hello.com/world?unit=testing",
			expected: "www.hello.com/world",
		},
		{
			name:     "test case 6",
			url:      "https://www.hello.com/world?unit=testing#foo",
			expected: "www.hello.com/world",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := normalizeURL(testCase.url)
			if err != nil {
				t.Errorf("%s failed, unexpected error: %v", testCase.name, err)
			}
			if result != testCase.expected {
				t.Errorf("%s failed, %s != %s", testCase.name, result, testCase.expected)
			}
		})
	}
}
