package storage

type Storage interface {
	Save(key string) error
	Read(key string) ([]byte, error)
}

func SaveTo(s Storage, key string) error {
	return s.Save(key)
}

func ReadFrom(s Storage, key string) ([]byte, error) {
	return s.Read(key)
}
