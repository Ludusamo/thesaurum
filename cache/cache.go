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

	HasNext() bool
	Next() Store
}

func StoreChain(s Store, topic string, data *Data) bool {
	succeeded := s.Store(topic, data)
	if s.HasNext() {
		return succeeded && StoreChain(s.Next(), topic, data)
	}
	return succeeded
}

func RetrieveChain(s Store, topic string) (*Data, bool) {
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

func ListChain(s Store) [][]string {
	lst := [][]string{s.List()}
	if s.HasNext() {
		lst = append(lst, ListChain(s.Next())...)
	}
	return lst
}
