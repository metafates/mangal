# Gache 

Gache is a dead simple file-based (or in-memory) cache library for Go with zero dependencies.

There are great caching libraries out there, but none of them fit my needs or not as simple as I'd like to.
So I decided to write my own and share it with the world. 🐳

> Work in progress, I wouldn't reccommend using it now

## Installation

```bash
go get github.com/metafates/gache
```

## Usage Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/metafates/gache"
	"net/http"
	"time"
)

type Pokemon struct {
	Height int
}

// Create new cache instance
var cache = gache.New[map[string]*Pokemon](&gache.Options{
	// Path to cache file
	// If not set, cache will be in-memory
	Path: ".cache/pokemons.json",

	// Lifetime of cache.
	// If not set, cache will never expire
	Lifetime: time.Hour,
})

// getPokemon will get a pokemon by name from API
// Gonna Cache Em' All!
func getPokemon(name string) (*Pokemon, error) {
	// check if Pokémon is in cache
	pokemons, expired, err := cache.Get()
	if err != nil {
		return nil, err
	}

	// if cache is expired, or Pokémon wasn't cached
	// Fetch it from API
	if pokemon, ok := pokemons[name]; !expired && ok {
		return pokemon, nil
	}

	// bla-bla-bla, boring stuff, etc...
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var pokemon Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, err
	}

	// okay, we got our Pokémon, let's cache it
	pokemons[name] = &pokemon
	_ = cache.Set(pokemons)

	return &pokemon, nil
}

func main() {
	start := time.Now()
	for i := 0; i < 3; i++ {
		_, _ = getPokemon("pikachu")
	}
	fmt.Println(time.Since(start))
}
```
