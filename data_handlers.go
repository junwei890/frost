package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/junwei890/rumbling/internal/database"
)

type crawlerConfig struct {
	db        *database.Queries
	links     map[string][]string
	domain    *url.URL
	mu        *sync.Mutex
	wg        *sync.WaitGroup
	control   chan struct{}
	maxVisits int
}

func (c *apiConfig) postData(w http.ResponseWriter, req *http.Request) {
	type reqData struct {
		Url string `json:"url"`
	}
	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		errorResponseWriter(w, http.StatusInternalServerError, err)
	}
	reqUrl := &reqData{}
	if err := json.Unmarshal(bytes, reqUrl); err != nil {
		errorResponseWriter(w, http.StatusBadRequest, err)
	}

	dom, err := url.Parse(reqUrl.Url)
	if err != nil {
		errorResponseWriter(w, http.StatusBadRequest, err)
	}
	crawler := &crawlerConfig{
		db:        c.db,
		links:     make(map[string][]string),
		domain:    dom,
		mu:        &sync.Mutex{},
		wg:        &sync.WaitGroup{},
		control:   make(chan struct{}, 5),
		maxVisits: 20,
	}

	crawler.initCrawl(reqUrl.Url)

	w.WriteHeader(http.StatusOK)
}
