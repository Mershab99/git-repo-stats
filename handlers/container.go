package handlers

import "github.com/coocood/freecache"

// Container will hold all dependencies for your application.
type Container struct {
	InMemCache *freecache.Cache
}

// NewContainer returns an empty or an initialized container for your handlers.
func NewContainer() (Container, error) {
	c := Container{
		InMemCache: freecache.NewCache(10 * 1024 * 1024),
	}
	return c, nil
}
