package nats

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/nats-io/nats.go"
)

type (
	JetStreamContext nats.JetStreamContext
	KVConfig   nats.KeyValueConfig
)

var (
	ErrBucketNotFound = nats.ErrBucketNotFound
)

type KV struct {
	kv nats.KeyValue
}

// UpsertKV will return or create a KeyValue store
func UpsertKV(js JetStreamContext, c *KVConfig) (*KV, error) {
	if c == nil || c.Bucket == "" {
		return nil, errors.New("invalid config")
	}
	kv, err := js.KeyValue(c.Bucket)
	switch {
	case errors.Is(err, ErrBucketNotFound):
		// some people may not like
		if kv, err = js.CreateKeyValue((*nats.KeyValueConfig)(c)); err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}
	return &KV{kv}, nil
}

// Put will place the new value for the key into the store.
func (kv *KV) Put(name string, value any) (int, error) {
	b, err := encode(value)
	if err != nil {
		return 0, err
	}
	n, err := kv.kv.Put(name, b)
	return int(n), err
}

// Get returns the latest value for the key.
func (kv *KV) Get(name string, dst any) (int, error) {
	entry, err := kv.kv.Get(name)
	if err != nil {
		return 0, err
	}
	err = decode(entry.Value(), dst)
	if err != nil {
		return 0, err
	}
	return int(entry.Revision()), err
}

func encode(value any) (b []byte, err error) {
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(value)
	if err != nil {
		return nil, nil
	}
	return buf.Bytes(), nil
}

func decode(p []byte, data any) (err error) {
	var buf = bytes.NewReader(p)
	err = gob.NewDecoder(buf).Decode(data)
	return
}
