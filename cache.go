package main

import (
	"fmt"
	"os"
)

type Store interface {
	Store(topic string, data string) bool
	Retrieve(topic string) (string, bool)
	Delete(topic string) bool
	List() []string

	HasNext() bool
	Next() Store
}

func StoreChain(s Store, topic string, data string) bool {
	succeeded := s.Store(topic, data)
	if s.HasNext() {
		return succeeded && StoreChain(s.Next(), topic, data)
	}
	return succeeded
}

func RetrieveChain(s Store, topic string) (string, bool) {
	data, found := s.Retrieve(topic)
	if !found && s.HasNext() {
		data, found = RetrieveChain(s.Next(), topic)
		if found {
			s.Store(topic, data)
		}
	}
	return data, found
}

func DeleteChain(s Store, topic string) bool {
	succeeded := s.Delete(topic)
	if s.HasNext() {
		return succeeded && DeleteChain(s.Next(), topic)
	}
	return succeeded
}

type InMemoryStore struct {
	data      map[string]string
	nextStore Store
}

func NewInMemoryStore(nextStore Store) *InMemoryStore {
	var s InMemoryStore
	s.data = make(map[string]string)
	s.nextStore = nextStore
	return &s
}

func (s *InMemoryStore) Store(topic string, data string) bool {
	fmt.Println("Storing in memory")
	s.data[topic] = data
	return true
}

func (s *InMemoryStore) Retrieve(topic string) (string, bool) {
	fmt.Println("Retrieving in memory")
	data, found := s.data[topic]
	return data, found
}

func (s *InMemoryStore) Delete(topic string) bool {
	delete(s.data, topic)
	return true
}

func (s *InMemoryStore) List() []string {
	topics := make([]string, len(s.data))
	i := 0
	for k := range s.data {
		topics[i] = k
		i++
	}
	return topics
}

func (s *InMemoryStore) Next() Store {
	return s.nextStore
}

func (s *InMemoryStore) HasNext() bool {
	return s.nextStore != nil
}

type FileStore struct {
	path      string
	nextStore Store
}

func (s *FileStore) getTopicPath(topic string) string {
	return fmt.Sprintf("%s/%s", s.path, topic)
}

func NewFileStore(path string, nextStore Store) *FileStore {
	var s FileStore
	s.path = path
	s.nextStore = nextStore
	return &s
}

func (s *FileStore) Store(topic string, data string) bool {
	fmt.Println("Storing in file")
	fout, err := os.Create(s.getTopicPath(topic))
	if err == nil {
		defer fout.Close()
		fout.WriteString(data)
		return true
	}
	return false
}

func (s *FileStore) Retrieve(topic string) (string, bool) {
	fmt.Println("Retrieving from file")
	data, err := os.ReadFile(s.getTopicPath(topic))
	if err != nil {
		return "", false
	}
	return string(data), true
}

func (s *FileStore) Delete(topic string) bool {
	err := os.Remove(s.getTopicPath(topic))
	return err == nil
}

func (s *FileStore) List() []string {
	files, err := os.Open(s.path)
	if err != nil {
		fmt.Println("error opening directory:", err)
		return nil
	}
	defer files.Close()

	fileInfos, err := files.ReadDir(-1)
	if err != nil {
		fmt.Println("error reading directory:", err)
		return nil
	}

	topics := make([]string, len(fileInfos))
	i := 0
	for _, info := range fileInfos {
		topics[i] = info.Name()
		i++
	}
	return topics
}

func (s *FileStore) Next() Store {
	return s.nextStore
}

func (s *FileStore) HasNext() bool {
	return s.nextStore != nil
}
