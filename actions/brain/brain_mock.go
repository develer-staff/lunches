package brain

import (
	"encoding/json"
	"errors"
)

type BrainMock map[string][]byte

func NewBrainMock() BrainMock {
	return make(BrainMock)
}

func (b BrainMock) Set(key string, val interface{}) error {
	encoded, err := json.Marshal(val)
	if err != nil {
		return err
	}

	b[key] = encoded

	return nil
}
func (b BrainMock) Read(key string) (string, error) {
	val, ok := b[key]

	if !ok {
		return "", errors.New("key not found")
	}

	return string(val), nil
}

func (b BrainMock) Get(key string, q interface{}) error {

	val, err := b.Read(key)

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), q)
}

func (b BrainMock) Close() error {
	b = nil
	return nil
}
