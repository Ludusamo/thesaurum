package cache

type Metadata struct {
	Size     int
	Datatype string // MIME Type
}

type Data struct {
	Meta Metadata
	Data []byte
}

type Store interface {
	Store(topic string, data *Data) bool
	Retrieve(topic string) (*Data, bool)
	Delete(topic string) bool
	List() []string
}

type Cache struct {
	layers []Store
}

func (cache *Cache) Add(store Store) *Cache {
	cache.layers = append(cache.layers, store)
	return cache
}

func (cache *Cache) Store(topic string, data *Data) bool {
	succeeded := true
	for _, store := range cache.layers {
		succeeded = succeeded && store.Store(topic, data)
	}
	return succeeded
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

func (cache *Cache) Delete(topic string) bool {
	succeeded := true
	for _, store := range cache.layers {
		succeeded = succeeded && store.Delete(topic)
	}
	return succeeded
}

func (cache *Cache) List() [][]string {
	var lst [][]string
	for _, store := range cache.layers {
		lst = append(lst, store.List())
	}
	return lst
}
