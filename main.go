package main

import (
	_ "embed"
	"encoding/json"
	"time"

	"golang.org/x/exp/rand"
)

var (
	//go:embed data/pokedex.json
	pokedexData []byte
)

func main() {
	rand.Seed(uint64(time.Now().UnixNano()))
	var pokemon []Pokemon
	json.Unmarshal(pokedexData, &pokemon)
	MainMenu(pokemon)
}

/*
type SquirdleMon struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Generation float64 `json:"gen"`
	Type1      string  `json:"type-1"`
	Type2      string  `json:"type-2"`
	Height     float64 `json:"Height"`
	Weight     float64 `json:"Weight"`

	description string
}

func ParseSquirdleDex() {
	b, err := os.ReadFile("squirdle-pokedex.json")
	utils.CheckError(err)

	om := orderedmap.New[string, []any]()
	err = json.Unmarshal(b, &om)
	utils.CheckError(err)

	pokemon := []SquirdleMon{}
	currentIdx := 1
	for pair := om.Oldest(); pair != nil; pair = pair.Next() {
		v := pair.Value

		if len(v) < 5 {
			panic("unexpected len")
		}

		gen, ok := v[0].(float64)
		utils.CheckOK(ok, "Gen wasn't an uint8")

		type1, ok := v[1].(string)
		utils.CheckOK(ok, "type1 wasn't a string")

		type2, ok := v[2].(string)
		utils.CheckOK(ok, "type2 wasn't a string")

		height, ok := v[3].(float64)
		utils.CheckOK(ok, "height wasn't a float32")

		weight, ok := v[4].(float64)
		utils.CheckOK(ok, "weight wasn't a float32")

		pokemon = append(pokemon, SquirdleMon{
			ID:         currentIdx,
			Name:       pair.Key,
			Generation: gen,
			Type1:      type1,
			Type2:      type2,
			Height:     height,
			Weight:     weight,
		})

		currentIdx++
	}

	d, err := json.Marshal(pokemon)
	utils.CheckError(err)
	utils.CheckError(os.WriteFile("clean-pokedex.json", d, 0664))
}

*/
