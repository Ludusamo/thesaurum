package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"thesaurum/cache"

	"github.com/julienschmidt/httprouter"
)

var chainCache cache.Cache

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/topic/", HandleList)
	router.GET("/topic/:topic", HandleGet)
	router.POST("/topic/:topic", HandlePost)

	dataPath := os.Getenv("DATA_FILEPATH")
	maxCache, err := strconv.Atoi(os.Getenv("MAX_MEMORY_CACHE"))
	if err != nil {
		log.Fatal("error parsing MAX_MEMORY_CACHE: ", err)
	}
	chainCache.Add(cache.NewInMemoryStore(maxCache)).Add(cache.NewFileStore(dataPath))

	log.Fatal(http.ListenAndServe(":5000", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Home")
}

func HandlePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	topic := p.ByName("topic")
	b, _ := io.ReadAll(r.Body)
	data := cache.Data{cache.Metadata{len(b), r.Header.Get("Content-Type")}, b}
	stored := chainCache.Store(topic, &data)
	if stored {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "success")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "failed to persist to all caches")
	}
}

func HandleGet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	topic := p.ByName("topic")
	data, found := chainCache.Retrieve(topic)
	if found {
		w.Header().Set("Content-Type", data.Meta.Datatype)
		w.WriteHeader(http.StatusOK)
		w.Write(data.Data)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "could not find data for topic")
	}
}

func HandleList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	topics := chainCache.List()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topics)
}
