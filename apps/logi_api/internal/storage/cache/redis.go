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

	// For production, you should load your CA certificate and potentially client certificates.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Set to false in production and provide a RootCA
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress:      []string{u.Host},
		Username:         u.User.Username(),
		Password:         password,
		TLSConfig:        tlsConfig,
		
	})
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
