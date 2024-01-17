package cache

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

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

func (s *FileStore) Store(topic string, data *Data) bool {
	fmt.Println("storing in file")
	fout, err := os.Create(s.getTopicPath(topic))
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer fout.Close()
	fmt.Fprintf(fout, "%d\n%s\n", data.Meta.Size, data.Meta.Datatype)
	fout.Write(data.Data)
	return true
}

func (s *FileStore) Retrieve(topic string) (*Data, bool) {
	fmt.Println("retrieving from file")
	fin, err := os.Open(s.getTopicPath(topic))
	if err != nil {
		return nil, false
	}

	scanner := bufio.NewScanner(bufio.NewReader(fin))
	if !scanner.Scan() {
		return nil, false
	}
	fileLen, err := strconv.Atoi(scanner.Text())
	if err != nil || !scanner.Scan() {
		return nil, false
	}
	dataType := scanner.Text()
	var data []byte
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}

	return &Data{Metadata{fileLen, dataType}, data}, true
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