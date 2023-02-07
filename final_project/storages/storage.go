package storages

import (
	"errors"

	"github.com/Antonov-Alexander/otus-go/final_project/types"
)

const (
	MemoryStorageType = iota
)

func GetStorage(storageType int) (types.Storage, error) {
	switch storageType {
	case MemoryStorageType:
		storage := &MemoryStorage{}
		storage.Init()
		return storage, nil
	}

	return nil, errors.New("undefined storage type")
}
