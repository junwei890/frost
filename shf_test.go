package main

import (
	"reflect"
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

func TestURLSFromHTML(t *testing.T) {
	testCases := []struct {
		name     string
		baseURL  string
		html     string
		expected []string
	}{
		{
			name:     "test case 1",
			baseURL:  "https://www.helloworld.com",
			html:     `<html><body><a href="https://www.hello.com"></a><a href = "/hello/world"></a></body></html>`,
			expected: []string{"https://www.hello.com", "https://www.helloworld.com/hello/world"},
		},
		{
			name:     "test case 2",
			baseURL:  "https://www.unittesting.com",
			html:     `<html><body><a href="/unit/testing/"></a><a href="https://www.npminstall.com?hello=world"></a></body></html>`,
			expected: []string{"https://www.unittesting.com/unit/testing/", "https://www.npminstall.com?hello=world"},
		},
		{
			name:     "test case 3",
			baseURL:  "https://www.neovim.com",
			html:     `<html><body><a href="/i/use?neovim=btw"></a><a href="/vim/mentioned"></a></body></html>`,
			expected: []string{"https://www.neovim.com/i/use?neovim=btw", "https://www.neovim.com/vim/mentioned"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := urlsFromHTML(testCase.html, testCase.baseURL)
			if err != nil {
				t.Errorf("%s failed, unexpected error: %v", testCase.name, err)
			}
			if !reflect.DeepEqual(result, testCase.expected) {
				t.Errorf("%s failed, %v != %v", testCase.name, result, testCase.expected)
			}
		})
	}
}
