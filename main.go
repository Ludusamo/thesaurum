package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"thesaurum/cache"

	"github.com/julienschmidt/httprouter"
)

var DataFilepath string
var store cache.Store

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/topic/", HandleList)
	router.GET("/topic/:topic", HandleGet)
	router.POST("/topic/:topic", HandlePost)

	DataFilepath = os.Getenv("DATA_FILEPATH")
	store = cache.NewInMemoryStore(cache.NewFileStore(DataFilepath, nil))

	log.Fatal(http.ListenAndServe(":5000", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Home")
}

func HandlePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	topic := p.ByName("topic")
	b, _ := io.ReadAll(r.Body)
	data := cache.Data{cache.Metadata{len(b), r.Header.Get("Content-Type")}, b}
	stored := cache.StoreChain(store, topic, &data)
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
	data, found := cache.RetrieveChain(store, topic)
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
	topics := cache.ListChain(store)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topics)
}
