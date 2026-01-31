package pokedex

import (
	"sync"

	"github.com/Nachsus/pokedexcli/internal/pokeapi"
)

type Pokedex struct {
	mu      sync.Mutex
	pokemon map[string]pokeapi.PokemonDetails
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		pokemon: make(map[string]pokeapi.PokemonDetails),
	}
}

func (p *Pokedex) Add(pokemon pokeapi.PokemonDetails) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pokemon[pokemon.Name] = pokemon
}

func (p *Pokedex) Get(name string) (pokeapi.PokemonDetails, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	pokemon, exists := p.pokemon[name]
	return pokemon, exists
}

func (p *Pokedex) GetAll() []pokeapi.PokemonDetails {
	p.mu.Lock()
	defer p.mu.Unlock()

	all := make([]pokeapi.PokemonDetails, 0, len(p.pokemon))
	for _, pokemon := range p.pokemon {
		all = append(all, pokemon)
	}
	return all
}

func (p *Pokedex) Has(name string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, exists := p.pokemon[name]
	return exists
}
