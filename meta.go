package filecache

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"time"
)

// MetaFromFile reads cache metadata from file
func MetaFromFile(path string) (meta *Meta, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	meta = &Meta{}
	if err = jsoniter.Unmarshal(data, &meta); err != nil {
		return
	}

	if meta.Key == "" {
		return nil, errors.New("cache meta file " + path + " is invalid: no key is defined")
	}
	if meta.Created == 0 {
		return nil, errors.New("cache meta file " + path + " is invalid: no created time is defined")
	}

	return
}

// Meta is a cache item metadata structure
// Short json keys are used to reduce file size
type Meta struct {
	// Key is a non-hashed unique for namespace item id
	// Required
	Key string `json:"k"`

	// Namespace is a cache item category folder name
	// If empty FileCache.NamespaceDefault will be set
	Namespace string `json:"n"`

	// OriginalName is an original file name; optional
	OriginalName string `json:"o,omitempty"`

	// TTL is a item's time-to-live value in seconds
	TTL int64 `json:"t"`

	// Created is a time when cache file was written
	// Do not set it by yourself
	Created int64 `json:"c"`

	// Fields is a map of any others metadata fields
	Fields map[string]interface{} `json:"f,omitempty"`
}

// IsExpired returns true if file is expired or if its TTL is 0
func (m *Meta) IsExpired() bool {
	if m.TTL == -1 {
		return false
	}
	now := time.Now().Unix()
	exp := m.Created + m.TTL
	return now > exp
}

// SaveToFile saves JSON-encoded metadata to file
func (m *Meta) SaveToFile(path string) error {
	if !pathIsMeta(path) {
		return errors.New(path + " is not a valid meta path: no '" + MetaPostfix + "' name postfix")
	}

	m.Created = time.Now().Unix()
	data, err := jsoniter.Marshal(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0744)
}
