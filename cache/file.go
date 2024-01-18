package cache

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type FileCache struct {
	path string
}

func (s *FileCache) getTopicPath(topic string) string {
	return fmt.Sprintf("%s/%s", s.path, topic)
}

func NewFileCache(path string) *FileCache {
	var s FileCache
	s.path = path
	return &s
}

func (s *FileCache) store(topic string, data *Data) error {
	log.Println("storing in file")
	fout, err := os.Create(s.getTopicPath(topic))
	if err != nil {
		log.Println(err)
		return err
	}
	defer fout.Close()
	fmt.Fprintf(fout, "%d\n%s\n", data.Meta.Size, data.Meta.Datatype)
	fout.Write(data.Data)
	return nil
}

func (s *FileCache) retrieve(topic string) (*Data, bool) {
	log.Println("retrieving from file")
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

func (s *FileCache) delete(topic string) error {
	err := os.Remove(s.getTopicPath(topic))
	return err
}

func (s *FileCache) list() []string {
	files, err := os.Open(s.path)
	if err != nil {
		log.Println("error opening directory:", err)
		return nil
	}
	defer files.Close()

	fileInfos, err := files.ReadDir(-1)
	if err != nil {
		log.Println("error reading directory:", err)
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
