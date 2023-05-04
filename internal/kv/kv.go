package kv

import (
	"errors"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	mgr nats.KeyValueManager
)

var (
	ErrKVUninitialized = errors.New("KV subsystem is uninitialized")
)

// InitializeKV must be called before anything else. It's safe to call multiple times.
func InitializeKV(js nats.JetStreamContext) {
	if mgr == nil {
		mgr = js
	}
}

// DefaultKVConfig returns a configuration with "mostly sane" defaults. Override
// with the following Option functions
func DefaultKVConfig(bucketName string) *nats.KeyValueConfig {
	return &nats.KeyValueConfig{
		Bucket: bucketName,
		// the zero-value for StorageType gives us file storage (as opposed to memory)
		// the other zero-values should yield a functional config
	}
}

// the intention here is to bury the kvOption type and only expose the functions
type kvOption func(c *nats.KeyValueConfig)

//nolint:revive // I know that kvOption is unexported, it's supposed to be
func WithTTL(d time.Duration) kvOption {
	return func(c *nats.KeyValueConfig) {
		c.TTL = d
	}
}

//nolint:revive // I know that kvOption is unexported, it's supposed to be
func WithReplicas(replicas int) kvOption {
	return func(c *nats.KeyValueConfig) {
		c.Replicas = replicas
	}
}

//nolint:revive // I know that kvOption is unexported, it's supposed to be
func WithDescription(desc string) kvOption {
	return func(c *nats.KeyValueConfig) {
		c.Description = desc
	}
}

// XXX: Not really sure we'd ever change this but...
//
//nolint:revive // I know that kvOption is unexported, it's supposed to be
func WithStorageType(st nats.StorageType) kvOption {
	return func(c *nats.KeyValueConfig) {
		c.Storage = st
	}
}

func CreateOrBindKVBucket(bucketName string, opts ...kvOption) (nats.KeyValue, error) {
	if mgr == nil {
		return nil, ErrKVUninitialized
	}
	kv, err := mgr.KeyValue(bucketName)
	if errors.Is(err, nats.ErrBucketNotFound) {
		cfg := DefaultKVConfig(bucketName)
		for _, o := range opts {
			o(cfg)
		}
		return mgr.CreateKeyValue(cfg)
	}
	return kv, err
}
