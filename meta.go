package filecache

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
