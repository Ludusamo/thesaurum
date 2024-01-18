package cache

import "errors"

type Metadata struct {
	Size     int
	Datatype string // MIME Type
}

type Data struct {
	Meta Metadata
	Data []byte
}

type CacheLayer interface {
	Store(topic string, data *Data) error
	Retrieve(topic string) (*Data, bool)
	Delete(topic string) error
	List() []string
}

type Cache struct {
	layers []CacheLayer
}

func (cache *Cache) Add(store CacheLayer) *Cache {
	cache.layers = append(cache.layers, store)
	return cache
}

func (cache *Cache) Store(topic string, data *Data) error {
	var err error
	for _, store := range cache.layers {
		err = errors.Join(err, store.Store(topic, data))
	}
	return err
}

func (cache *Cache) Retrieve(topic string) (*Data, bool) {
	for layer, store := range cache.layers {
		data, found := store.Retrieve(topic)
		if found {
			for i := 0; i < layer; i++ {
				cache.layers[i].Store(topic, data)
			}
			return data, found
		}
	}
	return nil, false
}

func (cache *Cache) Delete(topic string) error {
	var err error
	for _, store := range cache.layers {
		err = errors.Join(err, store.Delete(topic))
	}
	return err
}

func (cache *Cache) List() [][]string {
	var lst [][]string
	for _, store := range cache.layers {
		lst = append(lst, store.List())
	}
	return lst
}
