package main

import (
	"chonkle/utils"
	"errors"
	"fmt"
	"strconv"
)

const (
	SingleType = "Single-Type"
	NoneType   = "None"
)

type Pokemon struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Generation int     `json:"gen"`
	Type1      string  `json:"type-1"`
	Type2      string  `json:"type-2"`
	Height     float64 `json:"Height"`
	Weight     float64 `json:"Weight"`

	description string
}

func (p *Pokemon) SetDescription() {
	p.description = fmt.Sprintf("Gen %d, %s/%s, %sm, %skg",
		p.Generation,
		p.Type1,
		utils.CoalesceString(p.Type2, NoneType),
		strconv.FormatFloat(p.Height, 'f', 2, 64),
		strconv.FormatFloat(p.Weight, 'f', 2, 64))
}

func GetRandomPokemon(pokemon []Pokemon) Pokemon {
	return pokemon[utils.RandInt(0, len(pokemon)-1)]
}

func GetPokemonFromGenerationRange(minGen, maxGen int, pokemon []Pokemon) ([]Pokemon, error) {
	if minGen > maxGen {
		return nil, errors.New("minimum gen larger than max")
	}

	n := 0
	for _, p := range pokemon {
		if p.Generation >= minGen && p.Generation <= maxGen {
			pokemon[n] = p
			n++
		}
	}

	pokemon = pokemon[:n]
	return pokemon, nil
}
