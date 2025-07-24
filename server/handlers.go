package server

type CrawlerRes struct {
	URL string `json:"url"`
	Content []string `json:"content"`
}

