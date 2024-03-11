package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"thesaurum/cache"
)

var chainCache cache.Cache
var corsList string

func getEnv(env string, def string) string {
	if v, ok := os.LookupEnv(env); ok {
		return v
	}
	return def
}

func main() {
	corsList = getEnv("CORS_LIST", "")
	cacheLayers := strings.Split(getEnv("CACHE_LAYERS", "InMemory,File"), ",")
	for _, layer := range cacheLayers {
		if layer == "InMemory" {
			maxCache, err := strconv.Atoi(getEnv("MAX_MEMORY_CACHE", "1048576"))
			if err != nil {
				log.Fatal("error parsing MAX_MEMORY_CACHE: ", err)
			}
			chainCache.Add(cache.NewInMemoryCache(maxCache))
		} else if layer == "File" {
			dataPath := getEnv("DATA_FILEPATH", "data")
			chainCache.Add(cache.NewFileCache(dataPath))
		}
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /topic/", HandleList)
	router.HandleFunc("GET /topic/{topic}", HandleGet)
	router.HandleFunc("POST /topic/{topic}", HandlePost)
	router.HandleFunc("DELETE /topic/{topic}", HandleDelete)
	router.Handle("/editor/", http.StripPrefix("/editor/", http.FileServer(http.Dir("static"))))

	port, err := strconv.Atoi(getEnv("PORT", "5000"))
	if err != nil {
		log.Fatal("error parsing PORT: ", err)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

func enableCors(w *http.ResponseWriter) {
	if corsList != "" {
		(*w).Header().Set("Access-Control-Allow-Origin", corsList)
	}
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	topic := r.PathValue("topic")
	err := chainCache.Delete(topic)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "success")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	topic := r.PathValue("topic")
	b, _ := io.ReadAll(r.Body)
	data := cache.Data{
		Meta: cache.Metadata{
			Size: len(b),
			Datatype: r.Header.Get("Content-Type"),
		},
		Data: b,
	}
	err := chainCache.Store(topic, &data)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "success")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	topic := r.PathValue("topic")
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

func HandleList(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	topics := chainCache.List()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topics)
}
