package pokeapi

import (
	"time"

	"github.com/Nachsus/pokedexcli/internal/pokecache"
)

var cache = pokecache.NewCache(5 * time.Minute)
