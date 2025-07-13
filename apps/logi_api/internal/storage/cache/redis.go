package cache

import (
	"context"
	"crypto/tls"
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
