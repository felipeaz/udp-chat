package redis

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

// Cache implements Redis functions
type Cache struct {
	Host   string
	Port   string
	Expire string
}

// NewCache returns an instance of Cache
func NewCache(host, port, expire string) (*Cache, error) {
	cache := &Cache{
		Host:   host,
		Port:   port,
		Expire: expire,
	}

	_, err := cache.connect()
	return cache, err
}

// Connect initialize the cache.
func (c *Cache) connect() (redis.Conn, error) {
	h := fmt.Sprintf("%s:%s", c.Host, c.Port)

	conn, err := redis.Dial("tcp", h)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return conn, nil
}

// Set inserts a value into the Cache
func (c *Cache) Set(key string, value []byte) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, c.Expire)
	if err != nil {
		return err
	}

	return nil
}

// Get returns a value from Cache
func (c *Cache) Get(key string) ([]byte, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}

// Flush removes a value from Cache
func (c *Cache) Flush(key string) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("DEL", key)
	return err
}
