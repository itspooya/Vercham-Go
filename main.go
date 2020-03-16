package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"strings"
)

func WordCount(wordsToCount string) map[string]int {
	words := strings.Fields(wordsToCount)
	result := make(map[string]int)
	for _, word := range words {
		result[word] += 1
	}
	// By default in Go 1.12 and older maps are printed key-sorted order

	return result
}

type JsonRequest struct {
	Text string `json:"text"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check Header For json if => ok else => return Error As json
	if r.Header.Get("Content-Type") != "" {
		value := r.Header.Get("Content-Type")
		if value != "application/json" {
			msg := "Content-Type Header Should be set to application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return

		}
	} else if r.Header.Get("Content-Type") == "" {
		msg := "Content-Type Header Should be set to application/json"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	// Limit Max size to 1MB malicious or accident behavior
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var t JsonRequest
	if err := decoder.Decode(&t); err != nil {
		http.Error(w, "Invalid Field! Use format text:TextToFrequencyCheck", http.StatusBadRequest)
		return
	}
	Res := WordCount(t.Text)
	ResultJSON, err := json.Marshal(Res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(ResultJSON)
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	router := httprouter.New()
	router.POST("/", Index)
	log.Fatal(http.ListenAndServe(":8080", router))

}
