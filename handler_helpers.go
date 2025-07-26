package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func errorResponseWriter(w http.ResponseWriter, statusCode int, errMsg error) {
	log.Println(errMsg)

	type errRes struct {
		Error error `json:"error"`
	}
	res := errRes{
		Error: errMsg,
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(bytes); err != nil {
		log.Println(err)
	}
}
