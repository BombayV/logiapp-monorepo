package cache

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/valkey-io/valkey-go"
)

// Cache represents a connection to a cache like Redis.
type Cache struct {
	Client valkey.Client
}

// NewCache creates a new connection to the cache.
func NewCache(address string) (*Cache, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	password, _ := u.User.Password()

	// Build the client options
	opts := valkey.ClientOption{
		InitAddress: []string{u.Host},
		Username:    u.User.Username(),
		Password:    password,
	}

	// --- SOLUTION: Only enable TLS for "rediss://" scheme ---
	if u.Scheme == "rediss" {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // Should be false in production with a proper RootCA
		}
	}

	// Create the client with the appropriate options
	client, err := valkey.NewClient(opts)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the cache")
	return &Cache{Client: client}, nil
}

// Close closes the cache connection.
func (c *Cache) Close() {
	c.Client.Close()
}

// AddRevokedToken adds a token's JTI to the revocation list with an expiration.
func (c *Cache) AddRevokedToken(ctx context.Context, jti string, expiration time.Duration) error {
	return c.Client.Do(ctx, c.Client.B().Set().Key(jti).Value("revoked").Ex(expiration).Build()).Error()
}

// IsTokenRevoked checks if a token's JTI is in the revocation list.
func (c *Cache) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	exists, err := c.Client.Do(ctx, c.Client.B().Exists().Key(jti).Build()).AsInt64()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// Set stores a value in the cache with an expiration time
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.Client.Do(ctx, c.Client.B().Set().Key(key).Value(string(jsonValue)).Ex(expiration).Build()).Error()
}

// Get retrieves a value from the cache
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	result, err := c.Client.Do(ctx, c.Client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(result), dest)
}

// Delete removes a key from the cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.Client.Do(ctx, c.Client.B().Del().Key(key).Build()).Error()
}

// DeletePattern removes all keys matching a pattern
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := c.Client.Do(ctx, c.Client.B().Keys().Pattern(pattern).Build()).AsStrSlice()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return c.Client.Do(ctx, c.Client.B().Del().Key(keys...).Build()).Error()
}

// Exists checks if a key exists in the cache
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := c.Client.Do(ctx, c.Client.B().Exists().Key(key).Build()).AsInt64()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
