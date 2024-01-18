// Cache package implements a chained caching structure.
//
// The data that is stored in these caches are a uniform data struct with
// an arbitrary chunk of bytes tied with some additional metadata.
// Utility functions are exposed for propagating and storing this data across
// a chain of different cache layers.
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
	store(topic string, data *Data) error
	retrieve(topic string) (*Data, bool)
	delete(topic string) error
	list() []string
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
	for _, layer := range cache.layers {
		err = errors.Join(err, layer.store(topic, data))
	}
	return err
}

func (cache *Cache) Retrieve(topic string) (*Data, bool) {
	for layerNum, layer := range cache.layers {
		data, found := layer.retrieve(topic)
		if found {
			for i := 0; i < layerNum; i++ {
				cache.layers[i].store(topic, data)
			}
			return data, found
		}
	}
	return nil, false
}

func (cache *Cache) Delete(topic string) error {
	var err error
	for _, layer := range cache.layers {
		err = errors.Join(err, layer.delete(topic))
	}
	return err
}

func (cache *Cache) List() [][]string {
	var lst [][]string
	for _, layer := range cache.layers {
		lst = append(lst, layer.list())
	}
	return lst
}
