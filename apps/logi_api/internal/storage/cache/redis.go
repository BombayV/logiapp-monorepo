package cache

import (
	"crypto/tls"
	"log"
	"net/url"

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
