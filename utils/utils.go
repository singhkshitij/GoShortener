package utils

import (
	"net/url"

	uuid "github.com/satori/go.uuid"
	st "github.com/singhkshitij/GOShortener/store"
	"github.com/speps/go-hashids"
)

// Generator the type to generate keys(short urls)
type Generator func() string

// Factory is responsible to generate keys(short urls)
type Factory struct {
	store     st.Store
	generator Generator
}

// DefaultGenerator is the default url generator
func DefaultGenerator() string {
	id := uuid.NewV4()
	hd := hashids.NewData()
	hd.Salt = id.String()
	encoder, _ := hashids.NewWithData(hd)
	genID, _ := encoder.Encode([]int{3, 7, 9})
	return genID
}

// NewFactory receives a generator and a store and returns a new url Factory.
func NewFactory(generator Generator, store st.Store) *Factory {
	return &Factory{
		store:     store,
		generator: generator,
	}
}

// Gen generates the key.
func (factory *Factory) Gen(uri string) (key string, err error) {
	_, err = url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	}

	key = factory.generator()

	for {
		if v := factory.store.Get(key); v == "" {
			break
		}
		key = factory.generator()
	}
	return key, nil
}
