package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

var DataFilepath string
var store Store

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/topic/:topic", HandleGet)
	router.POST("/topic/:topic", HandlePost)

	DataFilepath = os.Getenv("DATA_FILEPATH")
	store = NewInMemoryStore(NewFileStore(DataFilepath, nil))

	log.Fatal(http.ListenAndServe(":5000", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Home")
}

func HandlePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	topic := p.ByName("topic")
	b, _ := io.ReadAll(r.Body)
	stored := StoreChain(store, topic, string(b))
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
	data, found := RetrieveChain(store, topic)
	if found {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, data)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "could not find data for topic")
	}
}
