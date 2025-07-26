package server

type CrawlerRes struct {
	URL string
	Doc string
}

type RakeRes struct {
	URL string
	Keywords []string
}
